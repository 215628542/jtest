package model

import (
	"time"
)

type Tt struct {
	Id         uint32    `gorm:"column:id"`
	TaskId     int64     `json:"task_id" gorm:"not null;type:tinyint(3);index:idx_task_id;index:idx_task_id_remark,priority:1;comment:任务id"` // 订单号
	Remark     string    `json:"remark" gorm:"not null;type:varchar(300);index:idx_remark;index:idx_task_id_remark,priority:2;comment:备注"`   // 订单号`
	CreateTime time.Time `json:"create_time" gorm:"index;not null;comment:创建时间"`
}

func (t Tt) Table() string {
	return "tt"
}
