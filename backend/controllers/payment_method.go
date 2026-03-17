package controllers

import (
	"account-management-system/backend/database"
	"account-management-system/backend/models"
	"account-management-system/backend/utils"

	"github.com/gin-gonic/gin"
)

// CreatePaymentMethodRequest 创建支付方式请求
type CreatePaymentMethodRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
}

// UpdatePaymentMethodRequest 更新支付方式请求
type UpdatePaymentMethodRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
}

// ListPaymentMethodsResponse 支付方式列表响应
type ListPaymentMethodsResponse struct {
	List     []PaymentMethodInfo `json:"list"`
	Total    int64               `json:"total"`
	Page     int                 `json:"page"`
	PageSize int                 `json:"page_size"`
}

// PaymentMethodInfo 支付方式信息
type PaymentMethodInfo struct {
	PaymentMethodID string `json:"payment_method_id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	SortOrder       int    `json:"sort_order"`
	CreateTime      string `json:"create_time"`
}

// CreatePaymentMethod 创建支付方式
func CreatePaymentMethod(c *gin.Context) {
	var req CreatePaymentMethodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, utils.ErrorResponse(2001, "参数错误"))
		return
	}

	// 检查支付方式名称是否已存在
	var count int64
	database.DB.Model(&models.PaymentMethod{}).Where("is_deleted = ?", 0).Where("name = ?", req.Name).Count(&count)
	if count > 0 {
		c.JSON(200, utils.ErrorResponse(2003, "支付方式名称已存在"))
		return
	}

	paymentMethod := models.PaymentMethod{
		PaymentMethodID: utils.GenerateID("pm"),
		Name:            req.Name,
		Description:     req.Description,
		SortOrder:       req.SortOrder,
		IsDeleted:       0,
	}

	if err := database.DB.Create(&paymentMethod).Error; err != nil {
		c.JSON(200, utils.ErrorResponse(5000, "创建支付方式失败"))
		return
	}

	// 记录操作日志
	LogOperation(c, "CREATE", "payment_method", "创建支付方式: "+req.Name)

	c.JSON(200, utils.SuccessResponse(gin.H{"payment_method_id": paymentMethod.PaymentMethodID}))
}

// ListPaymentMethods 获取支付方式列表
func ListPaymentMethods(c *gin.Context) {
	page := 1
	pageSize := 20
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

	query := database.DB.Model(&models.PaymentMethod{}).Where("is_deleted = ?", 0)

	// 模糊搜索
	if keyword != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	var total int64
	query.Count(&total)

	var paymentMethods []models.PaymentMethod
	offset := (page - 1) * pageSize
	query.Order("sort_order ASC").Offset(offset).Limit(pageSize).Find(&paymentMethods)

	// 转换为响应格式
	list := make([]PaymentMethodInfo, len(paymentMethods))
	for i, pm := range paymentMethods {
		list[i] = PaymentMethodInfo{
			PaymentMethodID: pm.PaymentMethodID,
			Name:            pm.Name,
			Description:     pm.Description,
			SortOrder:       pm.SortOrder,
			CreateTime:      pm.CreateTime.Format("2006-01-02 15:04:05"),
		}
	}

	c.JSON(200, utils.SuccessResponse(ListPaymentMethodsResponse{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}))
}

// UpdatePaymentMethod 更新支付方式
func UpdatePaymentMethod(c *gin.Context) {
	paymentMethodID := c.Param("id")
	if paymentMethodID == "" {
		c.JSON(200, utils.ErrorResponse(2001, "支付方式ID不能为空"))
		return
	}

	var req UpdatePaymentMethodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, utils.ErrorResponse(2001, "参数错误"))
		return
	}

	// 检查支付方式名称是否已被其他支付方式使用
	var count int64
	database.DB.Model(&models.PaymentMethod{}).Where("is_deleted = ?", 0).Where("name = ? AND payment_method_id != ?", req.Name, paymentMethodID).Count(&count)
	if count > 0 {
		c.JSON(200, utils.ErrorResponse(2003, "支付方式名称已存在"))
		return
	}

	updates := map[string]interface{}{
		"name":        req.Name,
		"description": req.Description,
		"sort_order":  req.SortOrder,
	}

	if err := database.DB.Model(&models.PaymentMethod{}).Where("is_deleted = ?", 0).Where("payment_method_id = ?", paymentMethodID).Updates(updates).Error; err != nil {
		c.JSON(200, utils.ErrorResponse(5000, "更新支付方式失败"))
		return
	}

	// 记录操作日志
	LogOperation(c, "UPDATE", "payment_method", "更新支付方式ID: "+paymentMethodID)

	c.JSON(200, utils.SuccessResponse(nil))
}

// DeletePaymentMethod 删除支付方式
func DeletePaymentMethod(c *gin.Context) {
	paymentMethodID := c.Param("id")
	if paymentMethodID == "" {
		c.JSON(200, utils.ErrorResponse(2001, "支付方式ID不能为空"))
		return
	}

	// 软删除
	if err := database.DB.Model(&models.PaymentMethod{}).Where("is_deleted = ?", 0).Where("payment_method_id = ?", paymentMethodID).Update("is_deleted", 1).Error; err != nil {
		c.JSON(200, utils.ErrorResponse(5000, "删除支付方式失败"))
		return
	}

	// 记录操作日志
	LogOperation(c, "DELETE", "payment_method", "删除支付方式ID: "+paymentMethodID)

	c.JSON(200, utils.SuccessResponse(nil))
}
