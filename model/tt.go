package model

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"time"
)

type Tt struct {
	Id         uint32    `gorm:"column:id"`
	TaskId     int64     `json:"task_id" gorm:"not null;type:tinyint(3);index:idx_task_id;index:idx_task_id_remark,priority:1;comment:任务id"`          // 订单号
	Remark     string    `json:"remark" gorm:"not null;type:varchar(300);default:'';index:idx_remark;index:idx_task_id_remark,priority:2;comment:备注"` // 订单号`
	CreateTime time.Time `json:"create_time" gorm:"index;not null;comment:创建时间"`
}

func (t Tt) Table() string {
	return "tt"
}

func (t *Tt) Get(ctx context.Context, db *gorm.DB, id int) (goodsIds []string, err error) {
	if t.Id < 1 {
		return goodsIds, errors.New("id is empty")
	}

	err = db.WithContext(ctx).Model(&Tt{}).Where("id = ?", id).Pluck("remark", &goodsIds).Error
	return
}
