package model

import (
	"time"
)

type Tt struct {
	Id         uint32    `gorm:"column:id"`
	TaskId     int64     `json:"task_id,omitempty"` // 订单号
	Remark     string    `json:"remark,omitempty"`
	CreateTime time.Time `json:"create_time" gorm:"index;not null;comment:创建时间"`
}

func (t Tt) Table() string {
	return "tt"
}
