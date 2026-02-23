package models

import (
	"time"

	"gorm.io/gorm"
)

// Department 部门模型
type Department struct {
	DepartmentID string         `gorm:"column:department_id;primaryKey;size:32" json:"department_id"`
	Name         string         `gorm:"column:name;uniqueIndex;size:50;not null" json:"name"`
	Description  string         `gorm:"column:description;size:200" json:"description"`
	SortOrder    int            `gorm:"column:sort_order;default:0" json:"sort_order"`
	CreateTime   time.Time      `gorm:"column:create_time;autoCreateTime" json:"create_time"`
	UpdateTime   time.Time      `gorm:"column:update_time;autoUpdateTime" json:"update_time"`
	DeleteTime   gorm.DeletedAt `gorm:"column:delete_time;index" json:"-"`
}

// TableName 指定表名
func (Department) TableName() string {
	return "t_department"
}
