package models

import (
	"time"
)

// Department 部门模型
type Department struct {
	DepartmentID string    `gorm:"column:department_id;primaryKey;size:32" json:"department_id"`
	Name         string    `gorm:"column:name;uniqueIndex;size:50;not null" json:"name"`
	Description  string    `gorm:"column:description;size:200" json:"description"`
	SortOrder    int       `gorm:"column:sort_order;default:0" json:"sort_order"`
	IsDeleted    int       `gorm:"column:is_deleted;default:0" json:"is_deleted"` // 0-未删除，1-已删除
	CreateTime   time.Time `gorm:"column:create_time;autoCreateTime" json:"create_time"`
	UpdateTime   time.Time `gorm:"column:update_time;autoUpdateTime" json:"update_time"`
}

// TableName 指定表名
func (Department) TableName() string {
	return "t_department"
}
