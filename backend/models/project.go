package models

import (
	"time"

	"gorm.io/gorm"
)

// Project 项目模型
type Project struct {
	ProjectID    string         `gorm:"column:project_id;primaryKey;size:32" json:"project_id"`
	Name         string         `gorm:"column:name;size:100;not null" json:"name"`
	Description  string         `gorm:"column:description;size:500" json:"description"`
	DepartmentID string         `gorm:"column:department_id;size:32" json:"department_id"` // 关联部门ID
	Status       int            `gorm:"column:status;default:1" json:"status"` // 1-进行中，0-已结束
	CreateTime   time.Time      `gorm:"column:create_time;autoCreateTime" json:"create_time"`
	UpdateTime   time.Time      `gorm:"column:update_time;autoUpdateTime" json:"update_time"`
	DeleteTime   gorm.DeletedAt `gorm:"column:delete_time;index" json:"-"`
}

// TableName 指定表名
func (Project) TableName() string {
	return "t_project"
}
