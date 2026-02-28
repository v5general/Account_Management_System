package controllers

import (
	"account-management-system/backend/database"
	"account-management-system/backend/models"
	"account-management-system/backend/utils"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username     string `json:"username" binding:"required"`
	Password     string `json:"password" binding:"required,min=8"`
	RealName     string `json:"real_name" binding:"required"`
	Role         string `json:"role" binding:"required,oneof=ADMIN EMPLOYEE FINANCE"`
	DepartmentID string `json:"department_id"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Username     string `json:"username" binding:"omitempty"`
	RealName     string `json:"real_name"`
	Role         string `json:"role" binding:"omitempty,oneof=ADMIN EMPLOYEE FINANCE"`
	DepartmentID string `json:"department_id"`
	Status       *int   `json:"status" binding:"omitempty,oneof=0 1"`
	Password     string `json:"password" binding:"omitempty,min=8"`
}

// UpdateMyAccountRequest 更新自己的账号请求
type UpdateMyAccountRequest struct {
	Username string `json:"username" binding:"omitempty"`
	RealName string `json:"real_name" binding:"required"`
	Password string `json:"password" binding:"omitempty,min=8"`
}

// ResetPasswordRequest 重置密码请求
type ResetPasswordRequest struct {
	Password string `json:"password" binding:"required,min=8"`
}

// ListUsersResponse 用户列表响应
type ListUsersResponse struct {
	List       []UserInfo `json:"list"`
	Total      int64      `json:"total"`
	Page       int        `json:"page"`
	PageSize   int        `json:"page_size"`
}

// CreateUser 创建用户
func CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, utils.ErrorResponse(2001, "参数错误"))
		return
	}

	// 验证用户名格式：只允许数字、字母、下划线
	if !isValidUsername(req.Username) {
		c.JSON(200, utils.ErrorResponse(2001, "用户名只能包含数字、字母和下划线"))
		return
	}

	// 检查用户名是否已存在
	var count int64
	database.DB.Model(&models.User{}).Where("is_deleted = ?", 0).Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		c.JSON(200, utils.ErrorResponse(2003, "用户名已存在"))
		return
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(200, utils.ErrorResponse(5000, "密码加密失败"))
		return
	}

	user := models.User{
		UserID:       utils.GenerateID("user"),
		Username:     req.Username,
		Password:     hashedPassword,
		RealName:     req.RealName,
		Role:         req.Role,
		DepartmentID: req.DepartmentID,
		Status:       1,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(200, utils.ErrorResponse(5000, "创建用户失败"))
		return
	}

	// 记录操作日志
	LogOperation(c, "CREATE", "user", "创建用户: "+req.Username)

	c.JSON(200, utils.SuccessResponse(gin.H{"user_id": user.UserID}))
}

// ListUsers 获取用户列表
func ListUsers(c *gin.Context) {
	page := 1
	pageSize := 20
	role := c.Query("role")
	departmentID := c.Query("department_id")
	keyword := c.Query("keyword")

	// 获取分页参数
	if p := c.Query("page"); p != "" {
		if val, err := parsePageParam(p); err == nil {
			page = val
		}
	}
	if ps := c.Query("page_size"); ps != "" {
		if val, err := parsePageParam(ps); err == nil {
			pageSize = val
		}
	}

	query := database.DB.Model(&models.User{}).Where("is_deleted = ?", 0)

	// 筛选条件
	if role != "" {
		query = query.Where("role = ?", role)
	}
	if departmentID != "" {
		query = query.Where("department_id = ?", departmentID)
	}
	if keyword != "" {
		query = query.Where("username LIKE ? OR real_name LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	var total int64
	query.Count(&total)

	// 查询用户列表（包含部门信息）
	var users []models.User
	offset := (page - 1) * pageSize
	query.Preload("Department", func(db *gorm.DB) *gorm.DB {
		return db.Where("is_deleted = ?", 0)
	}).Order("create_time DESC").Offset(offset).Limit(pageSize).Find(&users)

	// 转换为响应格式
	list := make([]UserInfo, len(users))
	for i, u := range users {
		deptName := ""
		if u.Department != nil {
			deptName = u.Department.Name
		}
		list[i] = UserInfo{
			UserID:         u.UserID,
			Username:       u.Username,
			RealName:       u.RealName,
			Role:           u.Role,
			DepartmentID:   u.DepartmentID,
			DepartmentName: deptName,
			Status:         u.Status,
			CreateTime:     u.CreateTime.Format("2006-01-02 15:04:05"),
		}
	}

	c.JSON(200, utils.SuccessResponse(ListUsersResponse{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}))
}

// UpdateUser 更新用户
func UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(200, utils.ErrorResponse(2001, "用户ID不能为空"))
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, utils.ErrorResponse(2001, "参数错误"))
		return
	}

	// 获取当前用户信息
	var currentUser models.User
	if err := database.DB.Where("user_id = ?", userID).First(&currentUser).Error; err != nil {
		c.JSON(200, utils.ErrorResponse(2002, "用户不存在"))
		return
	}

	// 如果修改用户名，检查是否重复
	if req.Username != "" && req.Username != currentUser.Username {
		var count int64
		database.DB.Model(&models.User{}).Where("is_deleted = ?", 0).Where("username = ?", req.Username).Count(&count)
		if count > 0 {
			c.JSON(200, utils.ErrorResponse(2003, "用户名已存在"))
			return
		}
	}

	updates := make(map[string]interface{})
	if req.Username != "" {
		updates["username"] = req.Username
	}
	if req.RealName != "" {
		updates["real_name"] = req.RealName
	}
	if req.Role != "" {
		updates["role"] = req.Role
	}
	if req.DepartmentID != "" {
		updates["department_id"] = req.DepartmentID
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	// 处理密码更新
	if req.Password != "" {
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			c.JSON(200, utils.ErrorResponse(5000, "密码加密失败"))
			return
		}
		updates["password"] = hashedPassword
	}

	if err := database.DB.Model(&models.User{}).Where("is_deleted = ?", 0).Where("user_id = ?", userID).Updates(updates).Error; err != nil {
		c.JSON(200, utils.ErrorResponse(5000, "更新用户失败"))
		return
	}

	// 记录操作日志
	LogOperation(c, "UPDATE", "user", "更新用户ID: "+userID)

	c.JSON(200, utils.SuccessResponse(nil))
}

// UpdateMyAccount 更新自己的账号信息
func UpdateMyAccount(c *gin.Context) {
	// 从上下文获取当前用户ID
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(200, utils.ErrorResponse(1002, "未登录"))
		return
	}

	var req UpdateMyAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, utils.ErrorResponse(2001, "参数错误"))
		return
	}

	// 获取当前用户信息
	var currentUser models.User
	if err := database.DB.Where("user_id = ?", userID).First(&currentUser).Error; err != nil {
		c.JSON(200, utils.ErrorResponse(2002, "用户不存在"))
		return
	}

	// 如果修改用户名，检查是否重复
	if req.Username != "" && req.Username != currentUser.Username {
		var count int64
		database.DB.Model(&models.User{}).Where("is_deleted = ?", 0).Where("username = ?", req.Username).Count(&count)
		if count > 0 {
			c.JSON(200, utils.ErrorResponse(2003, "用户名已存在"))
			return
		}
	}

	updates := make(map[string]interface{})
	if req.Username != "" {
		updates["username"] = req.Username
	}
	updates["real_name"] = req.RealName
	// 处理密码更新
	if req.Password != "" {
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			c.JSON(200, utils.ErrorResponse(5000, "密码加密失败"))
			return
		}
		updates["password"] = hashedPassword
	}

	if err := database.DB.Model(&models.User{}).Where("user_id = ?", userID).Updates(updates).Error; err != nil {
		c.JSON(200, utils.ErrorResponse(5000, "更新账号失败"))
		return
	}

	// 记录操作日志
	LogOperation(c, "UPDATE", "user", "更新账号信息")

	c.JSON(200, utils.SuccessResponse(nil))
}

// DeleteUser 删除用户
func DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(200, utils.ErrorResponse(2001, "用户ID不能为空"))
		return
	}

	if err := database.DB.Model(&models.User{}).Where("user_id = ?", userID).Update("is_deleted", 1).Error; err != nil {
		c.JSON(200, utils.ErrorResponse(5000, "删除用户失败"))
		return
	}

	c.JSON(200, utils.SuccessResponse(nil))
}

// ResetPassword 重置密码
func ResetPassword(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(200, utils.ErrorResponse(2001, "用户ID不能为空"))
		return
	}

	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, utils.ErrorResponse(2001, "参数错误"))
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(200, utils.ErrorResponse(5000, "密码加密失败"))
		return
	}

	if err := database.DB.Model(&models.User{}).Where("is_deleted = ?", 0).Where("user_id = ?", userID).Update("password", hashedPassword).Error; err != nil {
		c.JSON(200, utils.ErrorResponse(5000, "重置密码失败"))
		return
	}

	c.JSON(200, utils.SuccessResponse(nil))
}

func parsePageParam(s string) (int, error) {
	var val int
	if _, err := fmt.Sscanf(s, "%d", &val); err != nil {
		return 0, err
	}
	if val < 1 {
		return 1, nil
	}
	return val, nil
}

// isValidUsername 验证用户名格式
// 用户名只能包含数字、字母和下划线
func isValidUsername(username string) bool {
	if len(username) == 0 || len(username) > 50 {
		return false
	}
	for _, c := range username {
		if !isAlphaNumericOrUnderscore(c) {
			return false
		}
	}
	return true
}

// isAlphaNumericOrUnderscore 检查字符是否为字母、数字或下划线
func isAlphaNumericOrUnderscore(c rune) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}
