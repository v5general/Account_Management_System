package routes

import (
	"account-management-system/backend/controllers"
	"account-management-system/backend/middlewares"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 配置路由
func SetupRoutes(r *gin.Engine) {
	// 应用中间件
	r.Use(middlewares.CORSMiddleware())
	r.Use(middlewares.LoggerMiddleware())

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1
	v1 := r.Group("/api/v1")
	{
		// 认证接口（无需认证）
		auth := v1.Group("/auth")
		{
			auth.POST("/login", controllers.Login)
		}

		// 需要认证的接口
		authorized := v1.Group("")
		authorized.Use(middlewares.AuthMiddleware())
		{
			// 认证接口
			auth := authorized.Group("/auth")
			{
				auth.POST("/logout", controllers.Logout)
				auth.GET("/me", controllers.GetCurrentUser)
			}

			// 用户账号管理（用户可修改自己的账号信息）
			account := authorized.Group("/account")
			{
				account.PUT("/me", controllers.UpdateMyAccount)
			}

			// 员工列表（财务和管理员可访问，用于关联人员选择）
			users := authorized.Group("/users")
			{
				users.GET("", controllers.ListUsers)
			}

			// 用户管理（需要管理员权限）
			usersAdmin := authorized.Group("/users")
			usersAdmin.Use(middlewares.RequireAdmin())
			{
				usersAdmin.POST("", controllers.CreateUser)
				usersAdmin.PUT("/:id", controllers.UpdateUser)
				usersAdmin.DELETE("/:id", controllers.DeleteUser)
				usersAdmin.POST("/:id/reset-password", controllers.ResetPassword)
			}

			// 部门管理（需要管理员权限）
			departments := authorized.Group("/departments")
			departments.Use(middlewares.RequireAdmin())
			{
				departments.POST("", controllers.CreateDepartment)
				departments.GET("", controllers.ListDepartments)
				departments.GET("/:id", controllers.GetDepartment)
				departments.PUT("/:id", controllers.UpdateDepartment)
				departments.DELETE("/:id", controllers.DeleteDepartment)
				departments.GET("/:id/users", controllers.GetDepartmentUsers)
			}

			// 项目管理
			projects := authorized.Group("/projects")
			{
				projects.GET("", controllers.ListProjects)
				projects.GET("/:id", controllers.GetProject)
				// 创建、更新、删除需要财务权限
				projects.POST("", middlewares.RequireFinance(), controllers.CreateProject)
				projects.PUT("/:id", middlewares.RequireFinance(), controllers.UpdateProject)
				projects.DELETE("/:id", middlewares.RequireFinance(), controllers.DeleteProject)
			}

			// 费用分类（需要财务人员或管理员权限）
			categories := authorized.Group("/categories")
			{
				categories.GET("", controllers.ListCategories)
				// 创建、更新、删除需要财务权限
				categories.POST("", middlewares.RequireFinance(), controllers.CreateCategory)
				categories.PUT("/:id", middlewares.RequireFinance(), controllers.UpdateCategory)
				categories.DELETE("/:id", middlewares.RequireFinance(), controllers.DeleteCategory)
			}

			// 收支管理
			transactions := authorized.Group("/transactions")
			{
				// 创建和更新需要财务人员或管理员权限
				transactions.Use(middlewares.RequireFinance())
				transactions.POST("", controllers.CreateTransaction)
				transactions.PUT("/:id", controllers.UpdateTransaction)
				transactions.DELETE("/:id", controllers.DeleteTransaction)
				transactions.GET("/statistics", controllers.GetStatistics)
				// 收支审核（仅管理员可用）
				transactions.PUT("/:id/approve", middlewares.RequireAdmin(), controllers.ApproveTransaction)
				transactions.PUT("/:id/reject", middlewares.RequireAdmin(), controllers.RejectTransaction)
			}

			// 收支查询（所有认证用户可访问）
			transactions.GET("", controllers.ListTransactions)
			transactions.GET("/:id", controllers.GetTransaction)

			// 附件管理
			attachments := authorized.Group("/attachments")
			{
				attachments.POST("", controllers.UploadAttachment)
				attachments.GET("", controllers.ListAttachments)
				attachments.GET("/:id/download", controllers.DownloadAttachment)
				attachments.DELETE("/:id", controllers.DeleteAttachment)
			}

			// 报表导出
			reports := authorized.Group("/reports")
			reports.Use(middlewares.RequireFinance())
			{
				reports.GET("/:id/export", controllers.ExportReport)
			}
		}
	}
}
