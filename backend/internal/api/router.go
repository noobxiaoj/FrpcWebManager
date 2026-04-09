package api

import (
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xiaoj/frpc_webmanager/internal/api/handler"
	"github.com/xiaoj/frpc_webmanager/internal/middleware"
)

/**
 * 归一化文档语言参数。
 * 当前文档仅区分中文与英文两套版本，英文统一接收 en / en-US / en-GB 等前缀，
 * 其余情况全部回退为中文，保证接口行为稳定。
 *
 * @param lang 原始语言参数
 * @return 归一化后的语言标识，仅返回 zh 或 en
 */
func normalizeDocsLanguage(lang string) string {
	if strings.HasPrefix(strings.ToLower(strings.TrimSpace(lang)), "en") {
		return "en"
	}

	return "zh"
}

/**
 * 判断当前文件是否为英文文档变体。
 * 英文文档统一采用 `.en.md` 后缀命名，例如 `v0.9.2.en.md`。
 *
 * @param filename 文档文件名
 * @return 是否为英文文档文件
 */
func isEnglishDocFile(filename string) bool {
	return strings.HasSuffix(strings.ToLower(filename), ".en.md")
}

/**
 * 根据基础文件名与语言生成目标文档文件名。
 * 中文直接读取原文件；英文优先映射到 `.en.md` 文件。
 *
 * @param filename 基础文档文件名，例如 `v0.9.2.md`
 * @param lang 归一化后的语言标识
 * @return 对应语言的文档文件名
 */
func localizedDocFilename(filename string, lang string) string {
	if lang != "en" || !strings.HasSuffix(strings.ToLower(filename), ".md") {
		return filename
	}

	return strings.TrimSuffix(filename, ".md") + ".en.md"
}

// SetupRoutes 设置路由
func SetupRoutes(r *gin.Engine, taskHandler *handler.TaskHandler, serverHandler *handler.ServerHandler, settingsHandler *handler.SettingsHandler, ipWhitelistMiddleware *middleware.IPWhitelistMiddleware) {
	// 应用IP白名单中间件到所有API路由（除了健康检查）
	r.Use(ipWhitelistMiddleware.Middleware())

	// API 路由组
	api := r.Group("/api")
	{
		// 健康检查(输出日志以便在容器日志中查看健康检查状态)
		api.GET("/health", func(c *gin.Context) {
			c.String(http.StatusOK, "OK")
		})

		// 登录认证状态与登录行为
		auth := api.Group("/auth")
		{
			auth.GET("/status", settingsHandler.GetAuthStatus) // 获取当前认证状态
			auth.POST("/login", settingsHandler.Login)         // 登录并写入浏览器会话
			auth.POST("/logout", settingsHandler.Logout)       // 退出登录
		}

		// 任务管理路由
		tasks := api.Group("/tasks")
		{
			tasks.GET("", taskHandler.ListTasks)           // 获取任务列表
			tasks.POST("", taskHandler.CreateTask)         // 创建任务
			tasks.GET("/:id", taskHandler.GetTask)         // 获取任务详情
			tasks.PUT("/:id", taskHandler.UpdateTask)      // 更新任务
			tasks.DELETE("/:id", taskHandler.DeleteTask)   // 删除任务
			tasks.POST("/:id/start", taskHandler.StartTask) // 启动任务
			tasks.POST("/:id/stop", taskHandler.StopTask)   // 停止任务
			tasks.POST("/:id/reload", taskHandler.ReloadTask) // 重载任务
			tasks.GET("/:id/status", taskHandler.GetTaskStatus) // 获取任务状态
			tasks.GET("/:id/config", taskHandler.GetTaskConfig)   // 获取任务配置
		}

		// 服务器管理路由
		servers := api.Group("/servers")
		{
			servers.GET("", serverHandler.ListServers)                 // 获取服务器列表
			servers.POST("", serverHandler.CreateServer)               // 创建服务器
			servers.GET("/:id", serverHandler.GetServer)               // 获取服务器详情
			servers.PUT("/:id", serverHandler.UpdateServer)            // 更新服务器
			servers.POST("/:id/pause", serverHandler.PauseServer)      // 暂停服务器进程
			servers.POST("/:id/start", serverHandler.StartServer)      // 启动暂停中的服务器进程
			servers.POST("/:id/restart", serverHandler.RestartServer)  // 重启服务器进程
			servers.DELETE("/:id", serverHandler.DeleteServer)         // 删除服务器
			servers.PUT("/order", serverHandler.UpdateServerOrder)     // 更新服务器排序
			servers.PUT("/:id/lock", serverHandler.UpdateServerLock)   // 更新服务器锁定状态
			servers.GET("/:id/logs", serverHandler.GetLogs)            // 获取日志
			servers.POST("/:id/logs", serverHandler.AddLog)            // 添加日志
			servers.DELETE("/:id/logs", serverHandler.ClearLogs)       // 清空日志
			servers.GET("/:id/config", serverHandler.GetConfig)        // 获取配置文件路径
			servers.POST("/:id/config", serverHandler.GenerateConfig)  // 生成配置文件
		}

		// 系统设置路由
		settings := api.Group("/settings")
		{
			settings.GET("", settingsHandler.GetSettings)              // 获取设置
			settings.PUT("", settingsHandler.UpdateSettings)           // 更新设置
			settings.POST("/password", settingsHandler.CreatePassword) // 添加密码
			settings.PUT("/password", settingsHandler.UpdatePassword)  // 修改密码
			settings.DELETE("/password", settingsHandler.DeletePassword) // 删除密码
			settings.GET("/events", handler.HandleSettingsEvents)      // SSE 事件流
		}

		// 文档路由
		api.GET("/changelog/files", func(c *gin.Context) {
			docsPath := "./docs"
			lang := normalizeDocsLanguage(c.Query("lang"))

			// 读取目录中的所有文件
			files, err := os.ReadDir(docsPath)
			if err != nil {
				handler.ErrorResponse(c, handler.CodeInternalError, "无法读取文件列表", err)
				return
			}

			var versionFiles []string

			// 筛选出版本文件（统一返回基础文件名）。
			// 英文模式下只返回“存在英文翻译版本”的条目，避免前端切英文后仍混入中文原文。
			for _, file := range files {
				filename := file.Name()

				if file.IsDir() || !strings.HasSuffix(strings.ToLower(filename), ".md") {
					continue
				}

				if isEnglishDocFile(filename) {
					continue
				}

				if !strings.HasPrefix(strings.ToLower(filename), "v") {
					continue
				}

				if filename == "简介.md" {
					continue
				}

				if lang == "en" {
					englishFilename := localizedDocFilename(filename, lang)
					if _, statErr := os.Stat(filepath.Join(docsPath, englishFilename)); statErr != nil {
						continue
					}
				}

				versionFiles = append(versionFiles, filename)
			}

			// 按文件名排序
			sort.Strings(versionFiles)

			handler.SuccessResponse(c, gin.H{
				"files": versionFiles,
			})
		})

		// 获取文档内容
		api.GET("/changelog/file/:filename", func(c *gin.Context) {
			filename := c.Param("filename")
			docsPath := "./docs"
			lang := normalizeDocsLanguage(c.Query("lang"))

			// 安全检查:防止路径遍历攻击
			if strings.Contains(filename, "..") || strings.Contains(filename, "/") || strings.Contains(filename, "\\") {
				handler.ErrorResponse(c, handler.CodeBadRequest, "非法的文件名", nil)
				return
			}

			// 根据当前语言优先读取对应语言版本；
			// 若英文文档尚未提供，则自动回退到中文原始文件，避免关于页直接报错。
			targetFilename := filename
			if !isEnglishDocFile(filename) {
				localizedFilename := localizedDocFilename(filename, lang)
				localizedPath := filepath.Join(docsPath, localizedFilename)
				if _, statErr := os.Stat(localizedPath); statErr == nil {
					targetFilename = localizedFilename
				}
			}

			filePath := filepath.Join(docsPath, targetFilename)

			// 读取文件内容
			content, err := os.ReadFile(filePath)
			if err != nil {
				handler.ErrorResponse(c, handler.CodeNotFound, "文件不存在", err)
				return
			}

			handler.SuccessResponse(c, gin.H{
				"content": string(content),
			})
		})
	}
}
