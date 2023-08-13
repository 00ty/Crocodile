package models

import (
	"time"
)

type Kami struct {
	ID        uint
	Text      string
	Comment   string // 备注
	CreatedAt time.Time
	EndStr    string `gorm:"-"`
	ScUser    string // 生成用户
	State     int
	UseDate   time.Time // 使用时间
	Username  string    // 使用账户
	EndDate   time.Time // 到期时间
	Ext       string    // 扩展
}

// TableName 自定义表名
func (Kami) TableName() string {
	return "8eqq.n_kami"
}
