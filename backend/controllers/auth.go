package controllers

import (
	"account-management-system/backend/database"
	"account-management-system/backend/models"
	"account-management-system/backend/utils"

	"github.com/gin-gonic/gin"
)

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string      `json:"token"`
	User  *UserInfo   `json:"user"`
}

// UserInfo 用户信息
type UserInfo struct {
	UserID         string `json:"user_id"`
	Username       string `json:"username"`
	RealName       string `json:"real_name"`
	Role           string `json:"role"`
	DepartmentID   string `json:"department_id"`
	DepartmentName string `json:"department_name"` // 部门名称
	Status         int    `json:"status"`
	CreateTime     string `json:"create_time"`
}

// Login 用户登录
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, utils.ErrorResponse(2001, "参数错误"))
		return
	}

	var user models.User
	if err := database.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(200, utils.ErrorResponse(1001, "用户名或密码错误"))
		return
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		c.JSON(200, utils.ErrorResponse(1001, "用户名或密码错误"))
		return
	}

	// 检查用户状态
	if user.Status != 1 {
		c.JSON(200, utils.ErrorResponse(1003, "该账号已被禁用"))
		return
	}

	token, err := utils.GenerateToken(user.UserID, user.Username, user.Role)
	if err != nil {
		c.JSON(200, utils.ErrorResponse(5000, "生成令牌失败"))
		return
	}

	// 记录登录日志
	database.DB.Create(&models.OperationLog{
		LogID:         utils.GenerateID("log"),
		UserID:        user.UserID,
		OperationType: "LOGIN",
		Module:        "auth",
		Description:   "用户登录",
		RequestIP:     c.ClientIP(),
		Status:        1,
	})

	c.JSON(200, utils.SuccessResponse(LoginResponse{
		Token: token,
		User: &UserInfo{
			UserID:       user.UserID,
			Username:     user.Username,
			RealName:     user.RealName,
			Role:         user.Role,
			DepartmentID: user.DepartmentID,
			Status:       user.Status,
			CreateTime:   user.CreateTime.Format("2006-01-02 15:04:05"),
		},
	}))
}

// Logout 用户注销
func Logout(c *gin.Context) {
	// JWT是无状态的，注销主要在前端处理
	c.JSON(200, utils.SuccessResponse(nil))
}

// GetCurrentUser 获取当前用户信息
func GetCurrentUser(c *gin.Context) {
	userID := c.GetString("user_id")

	var user models.User
	if err := database.DB.Where("user_id = ?", userID).First(&user).Error; err != nil {
		c.JSON(200, utils.ErrorResponse(2002, "用户不存在"))
		return
	}

	c.JSON(200, utils.SuccessResponse(UserInfo{
		UserID:       user.UserID,
		Username:     user.Username,
		RealName:     user.RealName,
		Role:         user.Role,
		DepartmentID: user.DepartmentID,
		Status:       user.Status,
		CreateTime:   user.CreateTime.Format("2006-01-02 15:04:05"),
	}))
}
