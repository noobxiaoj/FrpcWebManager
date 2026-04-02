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

			// 读取目录中的所有文件
			files, err := os.ReadDir(docsPath)
			if err != nil {
				handler.ErrorResponse(c, handler.CodeInternalError, "无法读取文件列表", err)
				return
			}

			var versionFiles []string

			// 筛选出版本文件 (v开头的.md文件)
			for _, file := range files {
				filename := file.Name()
				if !file.IsDir() && strings.HasSuffix(filename, ".md") && strings.HasPrefix(strings.ToLower(filename), "v") {
					// 排除简介.md
					if filename != "简介.md" {
						versionFiles = append(versionFiles, filename)
					}
				}
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

			// 安全检查:防止路径遍历攻击
			if strings.Contains(filename, "..") || strings.Contains(filename, "/") || strings.Contains(filename, "\\") {
				handler.ErrorResponse(c, handler.CodeBadRequest, "非法的文件名", nil)
				return
			}

			filePath := filepath.Join(docsPath, filename)

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
