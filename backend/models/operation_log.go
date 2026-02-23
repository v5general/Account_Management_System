package models

import (
	"time"
)

// OperationLog 操作日志模型
type OperationLog struct {
	LogID          string    `gorm:"column:log_id;primaryKey;size:32" json:"log_id"`
	UserID         string    `gorm:"column:user_id;size:32;not null;index" json:"user_id"`
	OperationType  string    `gorm:"column:operation_type;size:50;not null" json:"operation_type"` // LOGIN/CREATE/UPDATE/DELETE/APPROVE
	Module         string    `gorm:"column:module;size:50" json:"module"`
	Description    string    `gorm:"column:description;size:500" json:"description"`
	RequestIP      string    `gorm:"column:request_ip;size:50" json:"request_ip"`
	Status         int       `gorm:"column:status;default:1" json:"status"` // 1-成功，0-失败
	CreateTime     time.Time `gorm:"column:create_time;autoCreateTime;index" json:"create_time"`

	// 关联
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName 指定表名
func (OperationLog) TableName() string {
	return "t_operation_log"
}
