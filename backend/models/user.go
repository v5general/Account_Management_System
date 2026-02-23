package models

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	UserID       string         `gorm:"column:user_id;primaryKey;size:32" json:"user_id"`
	Username     string         `gorm:"column:username;uniqueIndex;size:50;not null" json:"username"`
	Password     string         `gorm:"column:password;size:255;not null" json:"-"`
	RealName     string         `gorm:"column:real_name;size:50;not null" json:"real_name"` // 真实姓名
	Role         string         `gorm:"column:role;size:20;not null;default:EMPLOYEE" json:"role"` // ADMIN/EMPLOYEE/FINANCE
	DepartmentID string         `gorm:"column:department_id;size:32" json:"department_id"` // 关联部门ID
	Status       int            `gorm:"column:status;default:1" json:"status"` // 1-正常，0-禁用
	CreateTime   time.Time      `gorm:"column:create_time;autoCreateTime" json:"create_time"`
	UpdateTime   time.Time      `gorm:"column:update_time;autoUpdateTime" json:"update_time"`
	DeleteTime   gorm.DeletedAt `gorm:"column:delete_time;index" json:"-"`
}

// TableName 指定表名
func (User) TableName() string {
	return "t_user"
}
