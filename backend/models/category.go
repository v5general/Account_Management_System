package models

import (
	"time"
)

// Category 费用分类模型
type Category struct {
	CategoryID   string    `gorm:"column:category_id;primaryKey;size:32" json:"category_id"`
	Name         string    `gorm:"column:name;size:50;not null" json:"name"`
	Type         string    `gorm:"column:type;size:20;not null" json:"type"` // INCOME（收入）、EXPENSE（支出）
	Description  string    `gorm:"column:description;size:200" json:"description"`
	SortOrder    int       `gorm:"column:sort_order;default:0" json:"sort_order"`
	IsDeleted    int       `gorm:"column:is_deleted;default:0" json:"is_deleted"` // 0-未删除，1-已删除
	CreateTime   time.Time `gorm:"column:create_time;autoCreateTime" json:"create_time"`
	UpdateTime   time.Time `gorm:"column:update_time;autoUpdateTime" json:"update_time"`
}

// TableName 指定表名
func (Category) TableName() string {
	return "t_category"
}
