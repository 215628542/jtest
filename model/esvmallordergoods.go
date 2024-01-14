package model

import (
	"context"
	"errors"
)

type EsVmallOrderGoods struct {
	Id      uint64 `json:"id" gorm:"primary_key;autoIncrement:false"`
	OrderId uint64 `json:"order_id" gorm:"not null;index;comment:orders表的Id"`
	GoodsId string `json:"goods_id" gorm:"index;not null;type:varchar(50);comment:商品表的Id"`
}

func (m *EsVmallOrderGoods) GetById(ctx context.Context) error {
	if m.Id == 0 {
		return errors.New("id is empty")
	}
	return db.WithContext(ctx).Where("id=?", m.Id).Take(m).Error
}
