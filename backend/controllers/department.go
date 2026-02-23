package controllers

import (
	"account-management-system/backend/database"
	"account-management-system/backend/models"
	"account-management-system/backend/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateDepartmentRequest 创建部门请求
type CreateDepartmentRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
}

// UpdateDepartmentRequest 更新部门请求
type UpdateDepartmentRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
}

// CreateDepartment 创建部门
func CreateDepartment(c *gin.Context) {
	var req CreateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, utils.ErrorResponse(2001, "参数错误"))
		return
	}

	// 检查部门名称是否已存在
	var count int64
	database.DB.Model(&models.Department{}).Where("name = ?", req.Name).Count(&count)
	if count > 0 {
		c.JSON(200, utils.ErrorResponse(2003, "部门名称已存在"))
		return
	}

	department := models.Department{
		DepartmentID: utils.GenerateID("dept"),
		Name:         req.Name,
		Description:  req.Description,
		SortOrder:    req.SortOrder,
	}

	if err := database.DB.Create(&department).Error; err != nil {
		c.JSON(200, utils.ErrorResponse(5000, "创建部门失败"))
		return
	}

	// 记录操作日志
	LogOperation(c, "CREATE", "department", "创建部门: "+req.Name)

	c.JSON(200, utils.SuccessResponse(department))
}

// ListDepartments 获取部门列表
func ListDepartments(c *gin.Context) {
	var departments []models.Department
	if err := database.DB.Order("sort_order ASC").Find(&departments).Error; err != nil {
		c.JSON(200, utils.ErrorResponse(5000, "获取部门列表失败"))
		return
	}

	c.JSON(200, utils.SuccessResponse(departments))
}

// GetDepartment 获取部门详情
func GetDepartment(c *gin.Context) {
	id := c.Param("id")
	var department models.Department
	if err := database.DB.Where("department_id = ?", id).First(&department).Error; err != nil {
		c.JSON(200, utils.ErrorResponse(2002, "部门不存在"))
		return
	}

	c.JSON(200, utils.SuccessResponse(department))
}

// UpdateDepartment 更新部门
func UpdateDepartment(c *gin.Context) {
	id := c.Param("id")
	var req UpdateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, utils.ErrorResponse(2001, "参数错误"))
		return
	}

	var department models.Department
	if err := database.DB.Where("department_id = ?", id).First(&department).Error; err != nil {
		c.JSON(200, utils.ErrorResponse(2002, "部门不存在"))
		return
	}

	// 如果修改名称，检查是否重复
	if req.Name != "" && req.Name != department.Name {
		var count int64
		database.DB.Model(&models.Department{}).Where("name = ? AND department_id != ?", req.Name, id).Count(&count)
		if count > 0 {
			c.JSON(200, utils.ErrorResponse(2003, "部门名称已存在"))
			return
		}
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	updates["sort_order"] = req.SortOrder

	if err := database.DB.Model(&department).Updates(updates).Error; err != nil {
		c.JSON(200, utils.ErrorResponse(5000, "更新部门失败"))
		return
	}

	// 记录操作日志
	LogOperation(c, "UPDATE", "department", "更新部门: "+department.Name)

	c.JSON(200, utils.SuccessResponse(department))
}

// DeleteDepartment 删除部门
func DeleteDepartment(c *gin.Context) {
	id := c.Param("id")

	// 检查是否有用户关联
	var userCount int64
	database.DB.Model(&models.User{}).Where("department_id = ?", id).Count(&userCount)
	if userCount > 0 {
		c.JSON(200, utils.ErrorResponse(2001, "该部门下有用户，无法删除"))
		return
	}

	// 检查是否有项目关联
	var projectCount int64
	database.DB.Model(&models.Project{}).Where("department_id = ?", id).Count(&projectCount)
	if projectCount > 0 {
		c.JSON(200, utils.ErrorResponse(2001, "该部门下有关联项目，无法删除"))
		return
	}

	if err := database.DB.Where("department_id = ?", id).Delete(&models.Department{}).Error; err != nil {
		c.JSON(200, utils.ErrorResponse(5000, "删除部门失败"))
		return
	}

	// 记录操作日志
	LogOperation(c, "DELETE", "department", "删除部门ID: "+id)

	c.JSON(200, utils.SuccessResponse(nil))
}

// GetDepartmentUsers 获取部门下的用户列表
func GetDepartmentUsers(c *gin.Context) {
	id := c.Param("id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	var users []models.User
	var total int64

	database.DB.Model(&models.User{}).Where("department_id = ?", id).Count(&total)
	offset := (page - 1) * pageSize
	if err := database.DB.Where("department_id = ?", id).
		Offset(offset).
		Limit(pageSize).
		Find(&users).Error; err != nil {
		c.JSON(200, utils.ErrorResponse(5000, "获取用户列表失败"))
		return
	}

	c.JSON(200, utils.SuccessResponse(gin.H{
		"list":  users,
		"total": total,
	}))
}
