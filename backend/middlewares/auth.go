package middlewares

import (
	"net/http"
	"strings"

	"account-management-system/backend/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware JWT认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusOK, gin.H{
				"code":      1002,
				"message":   "令牌无效或过期",
				"data":      nil,
				"timestamp": utils.GetTimestamp(),
			})
			c.Abort()
			return
		}

		// Bearer token格式
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusOK, gin.H{
				"code":      1002,
				"message":   "令牌格式错误",
				"data":      nil,
				"timestamp": utils.GetTimestamp(),
			})
			c.Abort()
			return
		}

		token := parts[1]
		claims, err := utils.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":      1002,
				"message":   "令牌无效或过期",
				"data":      nil,
				"timestamp": utils.GetTimestamp(),
			})
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Next()
	}
}

// RequireAdmin 需要管理员权限
func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		if role != "ADMIN" {
			c.JSON(http.StatusOK, gin.H{
				"code":      1003,
				"message":   "权限不足",
				"data":      nil,
				"timestamp": utils.GetTimestamp(),
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequireFinance 需要财务人员或管理员权限
func RequireFinance() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		if role != "ADMIN" && role != "FINANCE" {
			c.JSON(http.StatusOK, gin.H{
				"code":      1003,
				"message":   "权限不足",
				"data":      nil,
				"timestamp": utils.GetTimestamp(),
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
