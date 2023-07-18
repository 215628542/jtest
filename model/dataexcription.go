package model

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type VmallDataExceptions struct {
	Id         uint32    `gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL"`
	Scene      string    `gorm:"column:scene;default:;comment:'异常场景'"`
	Source     string    `gorm:"column:source;default:;comment:'来源'"`
	RelationId string    `gorm:"column:relation_id;default:;comment:'关联id'"`
	OrderNum   string    `gorm:"column:order_num;default:;comment:'订单号'"`
	ReqData    string    `gorm:"column:req_data;default:;comment:'请求数据'"`
	ResData    string    `gorm:"column:res_data;default:;comment:'响应数据'"`
	TryNum     int8      `gorm:"column:try_num;default:0;comment:'重试次数'"`
	Status     string    `gorm:"column:status;default:NULL;comment:'场景状态'"`
	DelayTime  int32     `gorm:"column:delay_time;default:0;comment:'延迟处理时间，单位秒'"`
	ErrMsg     string    `gorm:"column:err_msg;default:;NOT NULL;comment:'错误信息'"`
	CreatedAt  time.Time `gorm:"column:created_at;default:NULL;comment:'创建时间'"`
	UpdatedAt  time.Time `gorm:"column:updated_at;default:NULL;comment:'更新时间'"`
}

func (v *VmallDataExceptions) TableName() string {
	return "data_exception"
}

func (v *VmallDataExceptions) BeforeCreate(*gorm.DB) error {
	v.CreatedAt = time.Now()
	return nil
}

func (v *VmallDataExceptions) BeforeUpdate(tx *gorm.DB) error {
	v.UpdatedAt = time.Now()
	return nil
}

func (v *VmallDataExceptions) Find(ctx context.Context, db *gorm.DB) []VmallDataExceptions {

	l := []VmallDataExceptions{}
	db.WithContext(ctx).Find(&l)
	return l
}
