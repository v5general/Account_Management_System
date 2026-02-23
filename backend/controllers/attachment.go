package controllers

import (
	"account-management-system/backend/database"
	"account-management-system/backend/models"
	"account-management-system/backend/utils"
	"os"

	"github.com/gin-gonic/gin"
)

// UploadResponse 上传响应
type UploadResponse struct {
	AttachmentID string `json:"attachment_id"`
	FileName     string `json:"file_name"`
	FileSize     int64  `json:"file_size"`
	FileType     string `json:"file_type"`
}

// UploadAttachment 上传附件
func UploadAttachment(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(200, utils.ErrorResponse(3001, "文件上传失败"))
		return
	}

	// 验证文件类型
	if !utils.ValidateFileType(file.Filename) {
		c.JSON(200, utils.ErrorResponse(3002, "文件格式不支持"))
		return
	}

	// 验证文件大小（100MB = 100 * 1024 * 1024 bytes）
	maxSize := int64(100 * 1024 * 1024)
	if file.Size > maxSize {
		c.JSON(200, utils.ErrorResponse(3003, "文件大小超限"))
		return
	}

	// 保存文件
	filePath, fileType, err := utils.UploadFile(file)
	if err != nil {
		c.JSON(200, utils.ErrorResponse(3001, "保存文件失败"))
		return
	}

	uploaderID := c.GetString("user_id")

	attachment := models.Attachment{
		AttachmentID: utils.GenerateID("att"),
		RecordID:     "", // 暂时为空，关联收支记录时更新
		FileName:     file.Filename,
		FilePath:     filePath,
		FileSize:     file.Size,
		FileType:     fileType,
		UploaderID:   uploaderID,
	}

	if err := database.DB.Create(&attachment).Error; err != nil {
		c.JSON(200, utils.ErrorResponse(5000, "保存附件记录失败"))
		return
	}

	c.JSON(200, utils.SuccessResponse(UploadResponse{
		AttachmentID: attachment.AttachmentID,
		FileName:     attachment.FileName,
		FileSize:     attachment.FileSize,
		FileType:     attachment.FileType,
	}))
}

// ListAttachments 获取附件列表
func ListAttachments(c *gin.Context) {
	recordID := c.Query("record_id")
	if recordID == "" {
		c.JSON(200, utils.ErrorResponse(2001, "记录ID不能为空"))
		return
	}

	var attachments []models.Attachment
	if err := database.DB.Where("record_id = ?", recordID).Find(&attachments).Error; err != nil {
		c.JSON(200, utils.ErrorResponse(5000, "查询附件失败"))
		return
	}

	c.JSON(200, utils.SuccessResponse(attachments))
}

// DownloadAttachment 下载附件
func DownloadAttachment(c *gin.Context) {
	attachmentID := c.Param("id")
	if attachmentID == "" {
		c.JSON(200, utils.ErrorResponse(2001, "附件ID不能为空"))
		return
	}

	var attachment models.Attachment
	if err := database.DB.Where("attachment_id = ?", attachmentID).First(&attachment).Error; err != nil {
		c.JSON(200, utils.ErrorResponse(2002, "附件不存在"))
		return
	}

	// 检查文件是否存在
	if _, err := os.Stat(attachment.FilePath); os.IsNotExist(err) {
		c.JSON(200, utils.ErrorResponse(2002, "文件不存在"))
		return
	}

	// 返回文件
	c.File(attachment.FilePath)
}

// DeleteAttachment 删除附件
func DeleteAttachment(c *gin.Context) {
	attachmentID := c.Param("id")
	if attachmentID == "" {
		c.JSON(200, utils.ErrorResponse(2001, "附件ID不能为空"))
		return
	}

	var attachment models.Attachment
	if err := database.DB.Where("attachment_id = ?", attachmentID).First(&attachment).Error; err != nil {
		c.JSON(200, utils.ErrorResponse(2002, "附件不存在"))
		return
	}

	// 权限检查
	role := c.GetString("role")
	userID := c.GetString("user_id")
	if role != "ADMIN" && role != "FINANCE" && attachment.UploaderID != userID {
		c.JSON(200, utils.ErrorResponse(1003, "权限不足"))
		return
	}

	// 已关联到收支记录的附件不可删除
	if attachment.RecordID != "" {
		c.JSON(200, utils.ErrorResponse(2001, "已关联到收支记录的附件不可删除"))
		return
	}

	// 删除文件
	if err := utils.DeleteFile(attachment.FilePath); err != nil {
		c.JSON(200, utils.ErrorResponse(5000, "删除文件失败"))
		return
	}

	// 删除数据库记录
	if err := database.DB.Delete(&attachment).Error; err != nil {
		c.JSON(200, utils.ErrorResponse(5000, "删除附件记录失败"))
		return
	}

	c.JSON(200, utils.SuccessResponse(nil))
}
