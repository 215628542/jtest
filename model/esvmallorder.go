package model

import (
	"context"
	"errors"
)

type EsVmallOrder struct {
	Id       uint64 `json:"id" gorm:"primary_key;autoIncrement:false"`
	Uid      string `json:"uid" gorm:"not null;type:varchar(50);index;comment:member表的Id"`
	OrderNum string `json:"order_num" gorm:"index;type:varchar(50);not null;comment:订单号"`
}

type EsData struct {
	EsVmallOrder
	VmallOrderGoods EsVmallOrderGoods `json:"vmall_order_goods"`
}

func (m *EsVmallOrder) GetById(ctx context.Context) error {
	if m.Id == 0 {
		return errors.New("id is empty")
	}
	return db.WithContext(ctx).Where("id=?", m.Id).Take(m).Error
}
