package model

import "time"

type VmallRuiyinxinPushLog struct {
	Id uint64 `json:"id" gorm:"primary_key;autoIncrement:false"`

	OrderId  uint64    `json:"order_id" gorm:"not null;index;comment:orders表的Id"`
	OrderNum string    `json:"order_num" gorm:"not null;type:varchar(50);comment:订单号"`
	Money    string    `json:"money" gorm:"not null;type:varchar(50);comment:金额"`
	Context  string    `json:"context" gorm:"not null;comment:备注"`
	PayAt    time.Time `json:"pay_at" gorm:"index;not null;comment:核销时间"`

	ThirdShopId  string `json:"third_shop_id" gorm:"not null;type:varchar(100);comment:三方的shop_id,有数据的时候核销需要推送信息"`
	ThirdOrderId string `json:"third_order_id" gorm:"not null;type:varchar(100);comment:瑞银信返回的orderId"`
	ThirdReqBody string `json:"third_req_body" gorm:"not null;type:text;comment:请求报文"`
	ThirdRspBody string `json:"third_rsp_body" gorm:"not null;type:text;comment:响应报文"`
	ThirdSysErr  string `json:"third_sys_err" gorm:"not null;comment:调用异常"`
	ThirdChannel string `json:"third_channel" gorm:"not null;type:varchar(300);default:'';comment:三方通道"`

	State  string `json:"state" gorm:"not null;type:varchar(20);comment:状态 SUCCESS-成功 FAIL-失败"`
	TryNum int    `json:"try_num" gorm:"not null;comment:尝试次数 最大5次"`

	ChannelId string `json:"channel_id" structs:"channel_id" gorm:"not null;index;type:varchar(50);comment:渠道ID"`

	CreatedAt time.Time `json:"created_at"  structs:"created_at"  gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at"  structs:"updated_at"  gorm:"not null"`
}

func (VmallRuiyinxinPushLog) TableName() string {
	return "vmall_ruiyinxin_push_logs"
}
