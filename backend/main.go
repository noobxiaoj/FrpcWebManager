package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xiaoj/frpc_webmanager/internal/api"
	"github.com/xiaoj/frpc_webmanager/internal/api/handler"
	"github.com/xiaoj/frpc_webmanager/internal/frpc"
	"github.com/xiaoj/frpc_webmanager/internal/middleware"
	"github.com/xiaoj/frpc_webmanager/internal/repository"
	"github.com/xiaoj/frpc_webmanager/internal/service"
)

// clearLogFiles 清空日志目录下所有日志文件的内容
func clearLogFiles(logDir string) error {
	// 检查日志目录是否存在
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		// 如果目录不存在,创建它
		return os.MkdirAll(logDir, 0755)
	}

	// 读取目录下的所有文件
	entries, err := os.ReadDir(logDir)
	if err != nil {
		return err
	}

	// 清空所有 .log 文件
	count := 0
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filename := entry.Name()
		if filepath.Ext(filename) == ".log" {
			logPath := filepath.Join(logDir, filename)
			// 以截断模式打开文件,清空内容
			file, err := os.OpenFile(logPath, os.O_TRUNC|os.O_WRONLY, 0644)
			if err != nil {
				log.Printf("清空日志文件 %s 失败: %v", filename, err)
				continue
			}
			defer file.Close()
			count++
		}
	}

	if count > 0 {
		log.Printf("已清空 %d 个日志文件", count)
	}

	return nil
}

func main() {
	// 创建数据目录
	// 优先使用环境变量，Docker 容器中使用 /app/data，开发环境使用 ../data
	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		// 检测是否在容器中运行（存在 /app/frontend/dist 目录）
		if _, err := os.Stat("/app/frontend/dist"); err == nil {
			dataDir = "/app/data"
		} else {
			dataDir = "../data"
		}
	}
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatalf("创建数据目录失败: %v", err)
	}

	// 清空日志文件（静默执行）
	logDir := filepath.Join(dataDir, "logs")
	_ = clearLogFiles(logDir) // 忽略错误

	// 初始化仓储
	store, err := repository.NewJSONStore(dataDir)
	if err != nil {
		log.Fatalf("初始化存储失败: %v", err)
	}

	taskRepo := repository.NewTaskRepository(store)
	processRepo := repository.NewProcessRepository(store)
	serverRepo := repository.NewServerRepository(store)
	settingsRepo, err := repository.NewSettingsRepository(dataDir)
	if err != nil {
		log.Fatalf("初始化设置仓库失败: %v", err)
	}

	// 初始化服务
	configService, err := service.NewConfigService(dataDir)
	if err != nil {
		log.Fatalf("初始化配置服务失败: %v", err)
	}

	// 初始化frpc管理器
	frpcManager, err := frpc.NewFrpcManager(filepath.Join(dataDir, "configs"))
	if err != nil {
		log.Fatalf("初始化frpc管理器失败: %v", err)
	}

	// 应用退出时清理所有frpc进程
	defer func() {
		log.Println("正在清理frpc进程...")
		processes := frpcManager.ListProcesses()
		for _, process := range processes {
			if err := frpcManager.StopServer(process.ServerAddr, process.ServerPort); err != nil {
				log.Printf("停止服务器 %s:%d 失败: %v", process.ServerAddr, process.ServerPort, err)
			}
		}
	}()

	// 先创建服务器服务(不注入 taskService,避免循环依赖)
	serverService := service.NewServerService(serverRepo, taskRepo, configService, frpcManager)

	// 创建设置服务
	settingsService := service.NewSettingsService(settingsRepo)

	// 创建IP白名单中间件
	ipWhitelistMiddleware := middleware.NewIPWhitelistMiddleware(settingsService)

	// 创建密码认证中间件
	passwordAuthMiddleware, err := middleware.NewPasswordAuthMiddleware(settingsService)
	if err != nil {
		log.Fatalf("初始化密码认证中间件失败: %v", err)
	}

	// 初始化任务服务(注入 serverService)
	taskService := service.NewTaskService(taskRepo, processRepo, frpcManager, serverService)

	// 恢复之前运行中的任务（静默执行）
	_ = taskService.RestoreRunningTasks() // 忽略错误，不影响启动

	// 初始化 API 处理器
	taskHandler := handler.NewTaskHandler(taskService)
	serverHandler := handler.NewServerHandler(serverService)
	settingsHandler := handler.NewSettingsHandler(settingsService, passwordAuthMiddleware)

	// 从 settings.json 读取前端端口配置
	frontendPort := 4500 // 默认端口
	settingsPath := filepath.Join(dataDir, "settings.json")
	if settingsData, err := os.ReadFile(settingsPath); err == nil {
		var settings map[string]interface{}
		if err := json.Unmarshal(settingsData, &settings); err == nil {
			if fp, ok := settings["frontendPort"].(float64); ok {
				frontendPort = int(fp)
			}
		}
	}

	// 环境变量可以覆盖配置文件中的端口
	if customPort := os.Getenv("PORT"); customPort != "" {
		var portInt int
		if err := json.Unmarshal([]byte(customPort), &portInt); err == nil {
			frontendPort = portInt
		}
	}

	// 设置 Gin 为 Release 模式，禁用调试日志
	gin.SetMode(gin.ReleaseMode)

	// 创建 Gin 路由（使用 gin.New() 避免默认日志中间件）
	// 生产环境禁用访问日志，避免日志污染
	r := gin.New()

	// 只添加 Recovery 中间件（捕获 panic），不添加 Logger 中间件
	r.Use(gin.Recovery())

	// CORS 中间件（容器内前后端同源，但仍保留以兼容开发环境）
	r.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			// 需要支持登录态 Cookie，因此有 Origin 时回显具体来源，避免 credentials 与 * 冲突。
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Vary", "Origin")
		} else {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		}
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// 密码访问控制中间件。
	// 放在路由注册前统一生效，这样所有 API 都能按配置自动拦截。
	r.Use(passwordAuthMiddleware.Middleware())

	// 健康检查日志中间件（只记录 /api/health 的访问，便于Docker健康检查调试）
	r.Use(func(c *gin.Context) {
		if c.Request.URL.Path == "/api/health" && c.Request.Method == "GET" {
			log.Printf("[健康检查] 收到检查请求，来源: %s", c.Request.RemoteAddr)
			c.Next()
			if c.Writer.Status() == 200 {
				log.Printf("[健康检查] 检查通过 ✓")
			} else {
				log.Printf("[健康检查] 检查失败 ✗ (状态码: %d)", c.Writer.Status())
			}
		} else {
			c.Next()
		}
	})

	// 设置路由
	api.SetupRoutes(r, taskHandler, serverHandler, settingsHandler, ipWhitelistMiddleware)

	// 静态文件服务 - 托管前端构建产物
	// 尝试多个可能的前端静态文件位置
	frontendDistPaths := []string{
		filepath.Join(dataDir, "../frontend/dist"), // 开发环境: ../data -> ../frontend/dist
		"/app/frontend/dist",                        // Docker 环境
		filepath.Join(filepath.Dir(dataDir), "frontend/dist"), // 备用路径
	}

	var frontendDist string
	for _, path := range frontendDistPaths {
		if _, err := os.Stat(path); err == nil {
			frontendDist = path
			break
		}
	}

	if frontendDist != "" {
		// SPA 路由处理 - 优先检查静态文件，不存在则返回 index.html
		r.NoRoute(func(c *gin.Context) {
			requestPath := c.Request.URL.Path

			// 如果是 API 请求但找不到路由，返回 404
			if strings.HasPrefix(requestPath, "/api/") {
				c.JSON(http.StatusNotFound, gin.H{"code": -1, "message": "API not found"})
				return
			}

			// 根路径直接返回 index.html
			if requestPath == "/" {
				indexPath := filepath.Join(frontendDist, "index.html")
				// log.Printf("[DEBUG] Serving root path with: %s", indexPath)
				c.File(indexPath)
				return
			}

			// 尝试直接提供静态文件
			// 使用 strings.TrimLeft 去掉开头的 /，然后用 filepath.Join
			filePath := filepath.Join(frontendDist, strings.TrimPrefix(requestPath, "/"))
			// log.Printf("[DEBUG] Request path: %s, File path: %s", requestPath, filePath)

			// 检查文件是否存在
			if _, err := os.Stat(filePath); err == nil {
				// log.Printf("[DEBUG] File exists, serving: %s", filePath)
				c.File(filePath)
				return
			}

			// log.Printf("[DEBUG] File not found, serving index.html")
			// 文件不存在，返回 index.html（SPA 路由）
			c.File(filepath.Join(frontendDist, "index.html"))
		})

		// 静默启用前端（不输出日志）
		// log.Printf("前端静态文件已启用，路径: %s", frontendDist)
	} else {
		log.Printf("警告: 前端构建产物不存在，跳过静态文件服务（已尝试路径: %v）", frontendDistPaths)
	}

	// 启动服务器
	port := fmt.Sprintf(":%d", frontendPort)
	// 启动服务器（不输出启动日志）
	if err := r.Run(port); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}
