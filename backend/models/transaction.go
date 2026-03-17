package models

import (
	"time"
)

// Transaction 收支流水模型
type Transaction struct {
	RecordID         string         `gorm:"column:record_id;primaryKey;size:32" json:"record_id"`
	Amount           float64        `gorm:"column:amount;type:decimal(15,2);not null" json:"amount"` // 正数=收入，负数=支出
	CategoryID       *string        `gorm:"column:category_id;size:32" json:"category_id"`
	ProjectID        *string        `gorm:"column:project_id;size:32" json:"project_id"` // 关联项目ID
	PersonID         *string        `gorm:"column:person_id;size:32" json:"person_id"`
	PaymentMethodID  *string        `gorm:"column:payment_method_id;size:32" json:"payment_method_id"` // 支付方式ID
	TransactionTime  time.Time      `gorm:"column:transaction_time;not null" json:"transaction_time"`
	Remark           string         `gorm:"column:remark;size:500" json:"remark"`
	Status           int            `gorm:"column:status;default:0" json:"status"` // 0-待审核，1-已审核，2-已驳回
	CreatorID        string         `gorm:"column:creator_id;size:32;not null" json:"creator_id"`
	IsDeleted        int            `gorm:"column:is_deleted;default:0" json:"is_deleted"` // 0-未删除，1-已删除
	CreateTime       time.Time      `gorm:"column:create_time;autoCreateTime" json:"create_time"`
	UpdateTime       time.Time      `gorm:"column:update_time;autoUpdateTime" json:"update_time"`

	// 关联
	Category      *Category      `gorm:"foreignKey:CategoryID;References:CategoryID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE" json:"category,omitempty"`
	Project       *Project       `gorm:"foreignKey:ProjectID;References:ProjectID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE" json:"project,omitempty"`
	Person        *User          `gorm:"foreignKey:PersonID;References:UserID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE" json:"person,omitempty"`
	PaymentMethod *PaymentMethod `gorm:"foreignKey:PaymentMethodID;References:PaymentMethodID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE" json:"payment_method,omitempty"`
	Creator       *User          `gorm:"foreignKey:CreatorID;References:UserID" json:"creator,omitempty"`
	Attachments   []Attachment   `gorm:"foreignKey:RecordID;References:RecordID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"attachments,omitempty"`
}

// TableName 指定表名
func (Transaction) TableName() string {
	return "t_transaction"
}
