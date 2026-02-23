package controllers

import (
	"account-management-system/backend/database"
	"account-management-system/backend/models"
	"account-management-system/backend/utils"
	"time"

	"github.com/gin-gonic/gin"
)

// CreateTransactionRequest 创建收支记录请求
type CreateTransactionRequest struct {
	Amount           float64  `json:"amount" binding:"required"`
	CategoryID       *string  `json:"category_id"`
	ProjectID        *string  `json:"project_id"`
	PersonID         *string  `json:"person_id"`
	TransactionTime  string   `json:"transaction_time" binding:"required"`
	Remark           string   `json:"remark"`
	AttachmentIDs    []string `json:"attachment_ids" binding:"required,min=1"`
}

// UpdateTransactionRequest 更新收支记录请求
type UpdateTransactionRequest struct {
	Amount           float64  `json:"amount" binding:"required"`
	CategoryID       *string  `json:"category_id"`
	ProjectID        *string  `json:"project_id"`
	PersonID         *string  `json:"person_id"`
	TransactionTime  string   `json:"transaction_time" binding:"required"`
	Remark           string   `json:"remark"`
	Status           int      `json:"status" binding:"omitempty,oneof=0 1 2"`
}

// ListTransactionsResponse 收支记录列表响应
type ListTransactionsResponse struct {
	List     []TransactionInfo `json:"list"`
	Total    int64             `json:"total"`
	Page     int               `json:"page"`
	PageSize int               `json:"page_size"`
}

// TransactionInfo 收支记录信息
type TransactionInfo struct {
	RecordID        string           `json:"record_id"`
	Amount          float64          `json:"amount"`
	CategoryID      *string          `json:"category_id"`
	ProjectID       *string          `json:"project_id"`
	Project         *ProjectSummary  `json:"project,omitempty"`
	PersonID        *string          `json:"person_id"`
	TransactionTime string           `json:"transaction_time"`
	Remark          string           `json:"remark"`
	Status          int              `json:"status"`
	CreatorID       string           `json:"creator_id"`
	CreateTime      string           `json:"create_time"`
	UpdateTime      string           `json:"update_time"`
	Category        *CategorySummary `json:"category,omitempty"`
	Person          *UserSummary     `json:"person,omitempty"`
	Creator         *UserSummary     `json:"creator,omitempty"`
	Attachments     []AttachmentInfo `json:"attachments,omitempty"`
}

// ProjectSummary 项目摘要
type ProjectSummary struct {
	ProjectID   string `json:"project_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// CategorySummary 分类摘要
type CategorySummary struct {
	CategoryID string `json:"category_id"`
	Name       string `json:"name"`
}

// UserSummary 用户摘要
type UserSummary struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
}

// AttachmentInfo 附件信息
type AttachmentInfo struct {
	AttachmentID string `json:"attachment_id"`
	FileName     string `json:"file_name"`
	FileSize     int64  `json:"file_size"`
	FileType     string `json:"file_type"`
}

// CreateTransaction 创建收支记录
func CreateTransaction(c *gin.Context) {
	var req CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, utils.ErrorResponse(2001, "参数错误"))
		return
	}

	// 验证支出记录必须有人员
	if req.Amount < 0 && req.PersonID == nil {
		c.JSON(200, utils.ErrorResponse(2001, "支出记录必须关联人员"))
		return
	}

	// 验证附件数量
	if len(req.AttachmentIDs) == 0 {
		c.JSON(200, utils.ErrorResponse(2001, "至少需要上传一个凭证附件"))
		return
	}

	// 解析交易时间
	transactionTime, err := time.Parse("2006-01-02 15:04:05", req.TransactionTime)
	if err != nil {
		c.JSON(200, utils.ErrorResponse(2001, "交易时间格式错误"))
		return
	}

	creatorID := c.GetString("user_id")

	transaction := models.Transaction{
		RecordID:        utils.GenerateID("record"),
		Amount:          req.Amount,
		CategoryID:      req.CategoryID,
		ProjectID:       req.ProjectID,
		PersonID:        req.PersonID,
		TransactionTime: transactionTime,
		Remark:          req.Remark,
		Status:          0, // 默认待审核
		CreatorID:       creatorID,
	}

	// 开始事务
	tx := database.DB.Begin()
	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		c.JSON(200, utils.ErrorResponse(5000, "创建收支记录失败"))
		return
	}

	// 关联附件
	if err := tx.Model(&models.Attachment{}).Where("attachment_id IN ?", req.AttachmentIDs).Updates(map[string]interface{}{
		"record_id": transaction.RecordID,
	}).Error; err != nil {
		tx.Rollback()
		c.JSON(200, utils.ErrorResponse(5000, "关联附件失败"))
		return
	}

	tx.Commit()

	c.JSON(200, utils.SuccessResponse(gin.H{"record_id": transaction.RecordID}))
}

// ListTransactions 获取收支记录列表
func ListTransactions(c *gin.Context) {
	page := 1
	pageSize := 20
	startTime := c.Query("start_time")
	endTime := c.Query("end_time")
	categoryID := c.Query("category_id")
	personID := c.Query("person_id")
	transType := c.DefaultQuery("type", "all")

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

	query := database.DB.Model(&models.Transaction{}).Preload("Category").Preload("Project").Preload("Person").Preload("Creator").Preload("Attachments")

	// 权限控制：员工只能查看本人关联的记录
	role := c.GetString("role")
	if role == "EMPLOYEE" {
		userID := c.GetString("user_id")
		query = query.Where("person_id = ?", userID)
	}

	// 筛选条件
	if startTime != "" {
		query = query.Where("transaction_time >= ?", startTime)
	}
	if endTime != "" {
		query = query.Where("transaction_time <= ?", endTime)
	}
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}
	projectID := c.Query("project_id")
	if projectID != "" {
		query = query.Where("project_id = ?", projectID)
	}
	if personID != "" {
		query = query.Where("person_id = ?", personID)
	}
	// 只显示已审核的记录（除非是管理员或财务人员）
	if role != "ADMIN" && role != "FINANCE" {
		query = query.Where("status = ?", 1)
	}
	if transType == "income" {
		query = query.Where("amount > 0")
	} else if transType == "expense" {
		query = query.Where("amount < 0")
	}

	var total int64
	query.Count(&total)

	var transactions []models.Transaction
	offset := (page - 1) * pageSize
	query.Order("transaction_time DESC").Offset(offset).Limit(pageSize).Find(&transactions)

	// 转换为响应格式
	list := make([]TransactionInfo, len(transactions))
	for i, t := range transactions {
		info := TransactionInfo{
			RecordID:        t.RecordID,
			Amount:          t.Amount,
			CategoryID:      t.CategoryID,
			ProjectID:       t.ProjectID,
			PersonID:        t.PersonID,
			TransactionTime: t.TransactionTime.Format("2006-01-02 15:04:05"),
			Remark:          t.Remark,
			Status:          t.Status,
			CreatorID:       t.CreatorID,
			CreateTime:      t.CreateTime.Format("2006-01-02 15:04:05"),
			UpdateTime:      t.UpdateTime.Format("2006-01-02 15:04:05"),
		}

		if t.Category != nil {
			info.Category = &CategorySummary{
				CategoryID: t.Category.CategoryID,
				Name:       t.Category.Name,
			}
		}
		if t.Project != nil {
			info.Project = &ProjectSummary{
				ProjectID:   t.Project.ProjectID,
				Name:        t.Project.Name,
				Description: t.Project.Description,
			}
		}
		if t.Person != nil {
			info.Person = &UserSummary{
				UserID:   t.Person.UserID,
				Username: t.Person.Username,
			}
		}
		if t.Creator != nil {
			info.Creator = &UserSummary{
				UserID:   t.Creator.UserID,
				Username: t.Creator.Username,
			}
		}

		// 附件信息
		attachments := make([]AttachmentInfo, len(t.Attachments))
		for j, a := range t.Attachments {
			attachments[j] = AttachmentInfo{
				AttachmentID: a.AttachmentID,
				FileName:     a.FileName,
				FileSize:     a.FileSize,
				FileType:     a.FileType,
			}
		}
		info.Attachments = attachments

		list[i] = info
	}

	c.JSON(200, utils.SuccessResponse(ListTransactionsResponse{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}))
}

// GetTransaction 获取收支记录详情
func GetTransaction(c *gin.Context) {
	recordID := c.Param("id")
	if recordID == "" {
		c.JSON(200, utils.ErrorResponse(2001, "记录ID不能为空"))
		return
	}

	var transaction models.Transaction
	if err := database.DB.Preload("Category").Preload("Project").Preload("Person").Preload("Creator").Preload("Attachments").
		Where("record_id = ?", recordID).First(&transaction).Error; err != nil {
		c.JSON(200, utils.ErrorResponse(2002, "记录不存在"))
		return
	}

	// 权限检查
	role := c.GetString("role")
	if role == "EMPLOYEE" {
		userID := c.GetString("user_id")
		if transaction.PersonID == nil || *transaction.PersonID != userID {
			c.JSON(200, utils.ErrorResponse(1003, "权限不足"))
			return
		}
	}

	c.JSON(200, utils.SuccessResponse(transaction))
}

// UpdateTransaction 更新收支记录
func UpdateTransaction(c *gin.Context) {
	recordID := c.Param("id")
	if recordID == "" {
		c.JSON(200, utils.ErrorResponse(2001, "记录ID不能为空"))
		return
	}

	var req UpdateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, utils.ErrorResponse(2001, "参数错误"))
		return
	}

	// 验证支出记录必须有人员
	if req.Amount < 0 && req.PersonID == nil {
		c.JSON(200, utils.ErrorResponse(2001, "支出记录必须关联人员"))
		return
	}

	// 解析交易时间
	transactionTime, err := time.Parse("2006-01-02 15:04:05", req.TransactionTime)
	if err != nil {
		c.JSON(200, utils.ErrorResponse(2001, "交易时间格式错误"))
		return
	}

	updates := map[string]interface{}{
		"amount":           req.Amount,
		"category_id":      req.CategoryID,
		"project_id":       req.ProjectID,
		"person_id":        req.PersonID,
		"transaction_time": transactionTime,
		"remark":           req.Remark,
	}

	if req.Status != 0 {
		updates["status"] = req.Status
	}

	if err := database.DB.Model(&models.Transaction{}).Where("record_id = ?", recordID).Updates(updates).Error; err != nil {
		c.JSON(200, utils.ErrorResponse(5000, "更新记录失败"))
		return
	}

	c.JSON(200, utils.SuccessResponse(nil))
}

// DeleteTransaction 删除收支记录（软删除）
func DeleteTransaction(c *gin.Context) {
	recordID := c.Param("id")
	if recordID == "" {
		c.JSON(200, utils.ErrorResponse(2001, "记录ID不能为空"))
		return
	}

	if err := database.DB.Delete(&models.Transaction{}, "record_id = ?", recordID).Error; err != nil {
		c.JSON(200, utils.ErrorResponse(5000, "删除记录失败"))
		return
	}

	// 记录操作日志
	LogOperation(c, "DELETE", "transaction", "删除收支记录: "+recordID)

	c.JSON(200, utils.SuccessResponse(nil))
}

// ApproveTransactionRequest 审核通过请求
type ApproveTransactionRequest struct {
	Remark string `json:"remark"`
}

// ApproveTransaction 审核通过收入记录
func ApproveTransaction(c *gin.Context) {
	recordID := c.Param("id")
	if recordID == "" {
		c.JSON(200, utils.ErrorResponse(2001, "记录ID不能为空"))
		return
	}

	var transaction models.Transaction
	if err := database.DB.Where("record_id = ?", recordID).First(&transaction).Error; err != nil {
		c.JSON(200, utils.ErrorResponse(2002, "记录不存在"))
		return
	}

	// 只有待审核状态的记录可以审核
	if transaction.Status != 0 {
		c.JSON(200, utils.ErrorResponse(2001, "该记录不是待审核状态"))
		return
	}

	// 更新状态为已审核
	if err := database.DB.Model(&transaction).Update("status", 1).Error; err != nil {
		c.JSON(200, utils.ErrorResponse(5000, "审核失败"))
		return
	}

	// 记录操作日志
	LogOperation(c, "APPROVE", "transaction", "审核通过收支记录: "+recordID)

	c.JSON(200, utils.SuccessResponse(nil))
}

// RejectTransactionRequest 驳回请求
type RejectTransactionRequest struct {
	Reason string `json:"reason" binding:"required"`
}

// RejectTransaction 驳回收入记录
func RejectTransaction(c *gin.Context) {
	recordID := c.Param("id")
	if recordID == "" {
		c.JSON(200, utils.ErrorResponse(2001, "记录ID不能为空"))
		return
	}

	var req RejectTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, utils.ErrorResponse(2001, "参数错误"))
		return
	}

	var transaction models.Transaction
	if err := database.DB.Where("record_id = ?", recordID).First(&transaction).Error; err != nil {
		c.JSON(200, utils.ErrorResponse(2002, "记录不存在"))
		return
	}

	// 只有待审核状态的记录可以驳回
	if transaction.Status != 0 {
		c.JSON(200, utils.ErrorResponse(2001, "该记录不是待审核状态"))
		return
	}

	// 更新状态为已驳回，并在备注中记录驳回原因
	updates := map[string]interface{}{
		"status": 2,
	}
	if req.Reason != "" {
		updates["remark"] = transaction.Remark + "\n[驳回原因: " + req.Reason + "]"
	}

	if err := database.DB.Model(&transaction).Updates(updates).Error; err != nil {
		c.JSON(200, utils.ErrorResponse(5000, "驳回失败"))
		return
	}

	// 记录操作日志
	LogOperation(c, "REJECT", "transaction", "驳回收支记录: "+recordID+", 原因: "+req.Reason)

	c.JSON(200, utils.SuccessResponse(nil))
}
