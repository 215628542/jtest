package model

import "gorm.io/gorm"

type WriteOffFail struct {
	gorm.Model
	OrderNum       string            `json:"order_num"  gorm:"column:order_num;varchar(50);default:0;comment:订单号"`
	TryNum         int               `json:"try_num"  gorm:"column:try_num;type:int;default:0;comment:重试次数"`
	WriteoffCode   string            `json:"writeoff_code" gorm:"not null;size:100;comment:核销码"`
	StaffId        string            `json:"staff_id" gorm:"not null;size:30;comment:员工id"`
	IsSuper        bool              `json:"is_super" gorm:" comment:是否超管"`
	OrgIds         stringArrayStruct `json:"orgIds" gorm:"type:json; comments:组织架构id"`
	IsBack         bool              `json:"is_back" gorm:" comment:是否后台操作"` // 是否后台操作，如果不是，则需要判断核销码
	WriteOffRemark string            `json:"write_off_remark" gorm:"not null;size:100;comment:核销备注"`
	ErrMsg         string            `json:"err_msg" gorm:"not null;size:100;comment:错误信息"`
	Status         string            `json:"status" gorm:"not null;size:100;comment:重推状态"`
}
