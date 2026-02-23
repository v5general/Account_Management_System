package controllers

import (
	"account-management-system/backend/database"
	"account-management-system/backend/models"
	"account-management-system/backend/utils"

	"github.com/gin-gonic/gin"
)

// CreateCategoryRequest 创建分类请求
type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Type        string `json:"type" binding:"required,oneof=INCOME EXPENSE"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
}

// UpdateCategoryRequest 更新分类请求
type UpdateCategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
}

// ListCategoriesResponse 分类列表响应
type ListCategoriesResponse struct {
	List     []CategoryInfo `json:"list"`
	Total    int64          `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"page_size"`
}

// CategoryInfo 分类信息
type CategoryInfo struct {
	CategoryID  string `json:"category_id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
	CreateTime  string `json:"create_time"`
}

// CreateCategory 创建费用分类
func CreateCategory(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, utils.ErrorResponse(2001, "参数错误"))
		return
	}

	// 检查分类名称是否已存在
	var count int64
	database.DB.Model(&models.Category{}).Where("is_deleted = ?", 0).Where("name = ?", req.Name).Count(&count)
	if count > 0 {
		c.JSON(200, utils.ErrorResponse(2003, "分类名称已存在"))
		return
	}

	category := models.Category{
		CategoryID:  utils.GenerateID("category"),
		Name:        req.Name,
		Type:        req.Type,
		Description: req.Description,
		SortOrder:   req.SortOrder,
		IsDeleted:   0,
	}

	if err := database.DB.Create(&category).Error; err != nil {
		c.JSON(200, utils.ErrorResponse(5000, "创建分类失败"))
		return
	}

	// 记录操作日志
	LogOperation(c, "CREATE", "category", "创建分类: "+req.Name)

	c.JSON(200, utils.SuccessResponse(gin.H{"category_id": category.CategoryID}))
}

// ListCategories 获取分类列表
func ListCategories(c *gin.Context) {
	page := 1
	pageSize := 20
	keyword := c.Query("keyword")
	categoryType := c.Query("type")

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

	query := database.DB.Model(&models.Category{}).Where("is_deleted = ?", 0)

	// 按类型筛选
	if categoryType != "" {
		query = query.Where("type = ?", categoryType)
	}

	// 模糊搜索
	if keyword != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	var total int64
	query.Count(&total)

	var categories []models.Category
	offset := (page - 1) * pageSize
	query.Order("sort_order ASC").Offset(offset).Limit(pageSize).Find(&categories)

	// 转换为响应格式
	list := make([]CategoryInfo, len(categories))
	for i, cat := range categories {
		list[i] = CategoryInfo{
			CategoryID:  cat.CategoryID,
			Name:        cat.Name,
			Type:        cat.Type,
			Description: cat.Description,
			SortOrder:   cat.SortOrder,
			CreateTime:  cat.CreateTime.Format("2006-01-02 15:04:05"),
		}
	}

	c.JSON(200, utils.SuccessResponse(ListCategoriesResponse{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}))
}

// UpdateCategory 更新分类
func UpdateCategory(c *gin.Context) {
	categoryID := c.Param("id")
	if categoryID == "" {
		c.JSON(200, utils.ErrorResponse(2001, "分类ID不能为空"))
		return
	}

	var req UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, utils.ErrorResponse(2001, "参数错误"))
		return
	}

	// 检查分类名称是否已被其他分类使用
	var count int64
	database.DB.Model(&models.Category{}).Where("is_deleted = ?", 0).Where("name = ? AND category_id != ?", req.Name, categoryID).Count(&count)
	if count > 0 {
		c.JSON(200, utils.ErrorResponse(2003, "分类名称已存在"))
		return
	}

	updates := map[string]interface{}{
		"name":        req.Name,
		"description": req.Description,
		"sort_order":  req.SortOrder,
	}

	if err := database.DB.Model(&models.Category{}).Where("is_deleted = ?", 0).Where("category_id = ?", categoryID).Updates(updates).Error; err != nil {
		c.JSON(200, utils.ErrorResponse(5000, "更新分类失败"))
		return
	}

	// 记录操作日志
	LogOperation(c, "UPDATE", "category", "更新分类ID: "+categoryID)

	c.JSON(200, utils.SuccessResponse(nil))
}

// DeleteCategory 删除分类
func DeleteCategory(c *gin.Context) {
	categoryID := c.Param("id")
	if categoryID == "" {
		c.JSON(200, utils.ErrorResponse(2001, "分类ID不能为空"))
		return
	}

	// 软删除
	if err := database.DB.Model(&models.Category{}).Where("is_deleted = ?", 0).Where("category_id = ?", categoryID).Update("is_deleted", 1).Error; err != nil {
		c.JSON(200, utils.ErrorResponse(5000, "删除分类失败"))
		return
	}

	// 记录操作日志
	LogOperation(c, "DELETE", "category", "删除分类ID: "+categoryID)

	c.JSON(200, utils.SuccessResponse(nil))
}
