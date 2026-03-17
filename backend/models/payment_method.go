package models

import (
	"time"
)

// PaymentMethod 支付方式模型
type PaymentMethod struct {
	PaymentMethodID string    `gorm:"column:payment_method_id;primaryKey;size:32" json:"payment_method_id"`
	Name            string    `gorm:"column:name;size:50;not null" json:"name"`
	Description     string    `gorm:"column:description;size:200" json:"description"`
	SortOrder       int       `gorm:"column:sort_order;default:0" json:"sort_order"`
	IsDeleted       int       `gorm:"column:is_deleted;default:0" json:"is_deleted"` // 0-未删除，1-已删除
	CreateTime      time.Time `gorm:"column:create_time;autoCreateTime" json:"create_time"`
	UpdateTime      time.Time `gorm:"column:update_time;autoUpdateTime" json:"update_time"`
}

// TableName 指定表名
func (PaymentMethod) TableName() string {
	return "t_payment_method"
}
