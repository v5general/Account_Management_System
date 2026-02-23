package models

import (
	"time"

	"gorm.io/gorm"
)

// Attachment 凭证附件模型
type Attachment struct {
	AttachmentID string         `gorm:"column:attachment_id;primaryKey;size:32" json:"attachment_id"`
	RecordID     string         `gorm:"column:record_id;size:32;not null;index" json:"record_id"`
	FileName     string         `gorm:"column:file_name;size:255;not null" json:"file_name"`
	FilePath     string         `gorm:"column:file_path;size:500;not null" json:"file_path"`
	FileSize     int64          `gorm:"column:file_size;not null" json:"file_size"`
	FileType     string         `gorm:"column:file_type;size:20;not null" json:"file_type"` // image/pdf
	UploadTime   time.Time      `gorm:"column:upload_time;autoCreateTime" json:"upload_time"`
	UploaderID   string         `gorm:"column:uploader_id;size:32;not null" json:"uploader_id"`
	DeleteTime   gorm.DeletedAt `gorm:"column:delete_time;index" json:"-"`

	// 关联
	Uploader *User `gorm:"foreignKey:UploaderID" json:"uploader,omitempty"`
}

// TableName 指定表名
func (Attachment) TableName() string {
	return "t_attachment"
}
