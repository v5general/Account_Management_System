package controllers

import (
	"account-management-system/backend/database"
	"account-management-system/backend/models"
	"account-management-system/backend/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateProjectRequest 创建项目请求
type CreateProjectRequest struct {
	Name         string `json:"name" binding:"required"`
	Description  string `json:"description"`
	DepartmentID string `json:"department_id"`
}

// UpdateProjectRequest 更新项目请求
type UpdateProjectRequest struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	DepartmentID string `json:"department_id"`
	Status       *int   `json:"status"`
}

// CreateProject 创建项目
func CreateProject(c *gin.Context) {
	var req CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, utils.ErrorResponse(2001, "参数错误"))
		return
	}

	// 检查项目名称是否已存在
	var count int64
	database.DB.Model(&models.Project{}).Where("is_deleted = ?", 0).Where("name = ?", req.Name).Count(&count)
	if count > 0 {
		c.JSON(200, utils.ErrorResponse(2003, "项目名称已存在"))
		return
	}

	project := models.Project{
		ProjectID:    utils.GenerateID("project"),
		Name:         req.Name,
		Description:  req.Description,
		DepartmentID: req.DepartmentID,
		Status:       1, // 默认进行中
	}

	if err := database.DB.Create(&project).Error; err != nil {
		c.JSON(200, utils.ErrorResponse(5000, "创建项目失败"))
		return
	}

	// 记录操作日志
	LogOperation(c, "CREATE", "project", "创建项目: "+req.Name)

	c.JSON(200, utils.SuccessResponse(project))
}

// ListProjects 获取项目列表
func ListProjects(c *gin.Context) {
	departmentID := c.Query("department_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	var projects []models.Project
	var total int64

	query := database.DB.Model(&models.Project{}).Where("is_deleted = ?", 0)
	if departmentID != "" {
		query = query.Where("department_id = ?", departmentID)
	}

	query.Count(&total)
	offset := (page - 1) * pageSize
	if err := query.Order("create_time DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&projects).Error; err != nil {
		c.JSON(200, utils.ErrorResponse(5000, "获取项目列表失败"))
		return
	}

	c.JSON(200, utils.SuccessResponse(gin.H{
		"list":  projects,
		"total": total,
	}))
}

// GetProject 获取项目详情
func GetProject(c *gin.Context) {
	id := c.Param("id")
	var project models.Project
	if err := database.DB.Where("is_deleted = ?", 0).Where("project_id = ?", id).First(&project).Error; err != nil {
		c.JSON(200, utils.ErrorResponse(2002, "项目不存在"))
		return
	}

	c.JSON(200, utils.SuccessResponse(project))
}

// UpdateProject 更新项目
func UpdateProject(c *gin.Context) {
	id := c.Param("id")
	var req UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, utils.ErrorResponse(2001, "参数错误"))
		return
	}

	var project models.Project
	if err := database.DB.Where("is_deleted = ?", 0).Where("project_id = ?", id).First(&project).Error; err != nil {
		c.JSON(200, utils.ErrorResponse(2002, "项目不存在"))
		return
	}

	// 如果修改名称，检查是否重复
	if req.Name != "" && req.Name != project.Name {
		var count int64
		database.DB.Model(&models.Project{}).Where("is_deleted = ?", 0).Where("name = ? AND project_id != ?", req.Name, id).Count(&count)
		if count > 0 {
			c.JSON(200, utils.ErrorResponse(2003, "项目名称已存在"))
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
	if req.DepartmentID != "" {
		updates["department_id"] = req.DepartmentID
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if err := database.DB.Model(&project).Updates(updates).Error; err != nil {
		c.JSON(200, utils.ErrorResponse(5000, "更新项目失败"))
		return
	}

	// 记录操作日志
	LogOperation(c, "UPDATE", "project", "更新项目: "+project.Name)

	c.JSON(200, utils.SuccessResponse(project))
}

// DeleteProject 删除项目
func DeleteProject(c *gin.Context) {
	id := c.Param("id")

	// 检查是否有收支记录关联
	var transCount int64
	database.DB.Model(&models.Transaction{}).Where("is_deleted = ?", 0).Where("project_id = ?", id).Count(&transCount)
	if transCount > 0 {
		c.JSON(200, utils.ErrorResponse(2001, "该项目有关联的收支记录，无法删除"))
		return
	}

	if err := database.DB.Model(&models.Project{}).Where("project_id = ?", id).Update("is_deleted", 1).Error; err != nil {
		c.JSON(200, utils.ErrorResponse(5000, "删除项目失败"))
		return
	}

	// 记录操作日志
	LogOperation(c, "DELETE", "project", "删除项目ID: "+id)

	c.JSON(200, utils.SuccessResponse(nil))
}
