package middlewares

import (
	"time"

	"account-management-system/backend/database"
	"account-management-system/backend/models"
	"account-management-system/backend/utils"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware 日志中间件
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		// 记录操作日志（仅记录POST/PUT/DELETE操作）
		if c.Request.Method != "GET" && c.Request.Method != "OPTIONS" {
			userID, exists := c.Get("user_id")
			if exists {
				log := &models.OperationLog{
					LogID:         utils.GenerateID("log"),
					UserID:        userID.(string),
					OperationType: getOperationType(c.Request.Method),
					Module:        getModuleFromPath(path),
					Description:   c.Request.Method + " " + path,
					RequestIP:     c.ClientIP(),
					Status:        1,
				}
				if c.Writer.Status() >= 400 {
					log.Status = 0
				}
				database.DB.Create(log)
			}
		}

		GinInfo(c, path, raw, latency)
	}
}

// GinInfo 输出请求日志
func GinInfo(c *gin.Context, path, raw string, latency time.Duration) {
	// 可以接入zap等日志库
	// 这里简单输出
}

// getOperationType 获取操作类型
func getOperationType(method string) string {
	switch method {
	case "POST":
		return "CREATE"
	case "PUT", "PATCH":
		return "UPDATE"
	case "DELETE":
		return "DELETE"
	default:
		return "QUERY"
	}
}

// getModuleFromPath 从路径获取模块名
func getModuleFromPath(path string) string {
	if len(path) > 5 {
		return path[5:] // 去掉 /api 前缀
	}
	return path
}
