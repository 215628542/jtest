package model

import (
	"gorm.io/gorm"
	"time"
)

type VmallOrderPay struct {
	Id uint64 `json:"id" gorm:"primary_key;autoIncrement:false"`

	OrderId       uint64       `json:"order_id" gorm:"not null;index;comment:orders表的Id"`
	TransactionId string       `json:"transaction_id" gorm:"index;type:varchar(50);not null;comment:微信支付单号"`
	TradState     string       `json:"trade_state" gorm:"not null;type:varchar(50);comment:交易状态"`
	SuccessedAt   time.Time    `json:"successed_at" gorm:"not null"`
	Openid        string       `json:"openid" gorm:"not null;type:varchar(50);comment:支付者信息"`
	Total         uint64       `json:"total" gorm:"not null;comment:总金额"`
	PayerTotal    uint64       `json:"payer_total" gorm:"not null;comment:用户支付金额"`
	PayOrderNum   string       `json:"pay_order_num" gorm:"not null;type:varchar(50);comment:支付微服务订单号"`
	DiscountMoney uint64       `json:"discount_money" gorm:"not null;comment:优惠金额"`
	ActuallyMoney uint64       `json:"actually_money" gorm:"not null;comment:实付金额"`
	DiscountInfo  string       `json:"discount_info" gorm:"not null;type:text;comment:优惠信息"`
	AttchInfo     string       `json:"attch_info" gorm:"not null;type:varchar(200);comment:attach信息"`
	BankType      string       `json:"card_type" gorm:"not null;size:50;comment:支付卡类型"`
	UseCouponIds  GormStrSlice `json:"use_coupon_ids" gorm:"type:json;comment:使用的优惠券id"`
	PayType       string       `json:"pay_type" gorm:"not null;size:50;default:'wechatPay';comment:支付类型"`
	ReservedField string       `json:"reserved_field" gorm:"not null;size:300;comment:预留字段，数币支付代表action"`
	CreatedAt     time.Time    `json:"created_at"  structs:"created_at"  gorm:"not null"`
	UpdatedAt     time.Time    `json:"updated_at"  structs:"updated_at"  gorm:"not null"`
	TestJson      GormJson     `json:"test_json" gorm:"type:json;comment:测试json类型"`
	TestSlice     GormStrSlice `json:"test_slice" gorm:"type:json;comment:测试slice类型"`
	TestBool      bool         `json:"test_bool" gorm:"type:json;comment:测试bool类型"`
}

func (m *VmallOrderPay) BeforeCreate(*gorm.DB) error {
	if m.Id == 0 {
		m.Id = uint64(2)
	}
	return nil
}
