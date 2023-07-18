package testmodel

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type Goods struct {
	Id                  uint64 `json:"id,omitempty"`
	CreatedAt           string `json:"created_at,omitempty"`
	UpdatedAt           string `json:"updated_at,omitempty"`
	DeletedAt           string `json:"deleted_at,omitempty"`
	Name                string `json:"name,omitempty"`              // 商品名称
	OrgId               uint64 `json:"org_id,omitempty"`            // 组织ID
	GoodsType           uint8  `json:"goods_type,omitempty"`        // 商品类型，单选（1.卡密权益/2.直充权益/3.微信代金券/4.自提 5.邮寄）
	SupplyType          uint8  `json:"supply_type,omitempty"`       // 供应类型，单选（下拉框（1.商城零售/2.直充频道）
	EquityId            uint64 `json:"equity_id,omitempty"`         // 权益ID
	StockNumber         int16  `json:"stock_number,omitempty"`      // 目前库存
	InitStockNumber     int16  `json:"init_stock_number,omitempty"` // 初始化库存
	Images              string `json:"images,omitempty"`
	SalesType           uint8  `json:"sales_type,omitempty"`        // 销售类型，单选 1.无门槛商品/2.兑换商品/3.现金商品/4.接口积分兑换商品
	ExchangeNumber      uint64 `json:"exchange_number,omitempty"`   // 兑换值，本地积分
	SalesPrice          string `json:"sales_price,omitempty"`       // 销售价
	MarketPrice         string `json:"market_price,omitempty"`      // 市场价
	MemberBuyType       uint8  `json:"member_buy_type,omitempty"`   // 仅限会员购买，1.不限 2.实名会员
	BuyLimit            uint32 `json:"buy_limit,omitempty"`         // 限购数量
	IsExpiredRefund     int8   `json:"is_expired_refund,omitempty"` // 过期是否自动退款
	Tags                string `json:"tags,omitempty"`              // 标签
	Sort                uint32 `json:"sort,omitempty"`              // 排序
	IsSale              int8   `json:"is_sale,omitempty"`           // 是否上架 1.是 0.否
	ChannelId           string `json:"channel_id,omitempty"`
	AppId               string `json:"app_id,omitempty"`
	CustomerAppId       string `json:"customer_app_id,omitempty"`
	Detail              string `json:"detail,omitempty"`               // 详情-富文本
	UseDetail           string `json:"use_detail,omitempty"`           // 使用说明-富文本
	SaleStockNumber     uint16 `json:"sale_stock_number,omitempty"`    // 总库存
	MemberType          string `json:"member_type,omitempty"`          // 会员类型
	ApiExchangeNumber   uint64 `json:"api_exchange_number,omitempty"`  // 接口积分兑换值
	IsHidden            int8   `json:"is_hidden,omitempty"`            // 是否在前端页面隐藏 1.是 0.否
	OffshelfTime        string `json:"offshelf_time,omitempty"`        // 留空不自动下架，下架时间
	IsDistribution      int8   `json:"is_distribution,omitempty"`      // 是否分销商品，1-是 0-否
	LimitCardTryNum     uint64 `json:"limit_card_try_num,omitempty"`   // 可尝试次数
	LimitCard           string `json:"limit_card,omitempty"`           // 限制的银行卡
	LimitCardTip        string `json:"limit_card_tip,omitempty"`       // 限卡提示文本
	IsLimitCard         int8   `json:"is_limit_card,omitempty"`        // 是否限卡商品 1-是 0-否
	IsDistributionTask  int8   `json:"is_distribution_task,omitempty"` // 是否分销任务商品，1-是 0-否
	AdBg                string `json:"ad_bg,omitempty"`
	RelationLink        string `json:"relation_link,omitempty"`
	ShareBg             string `json:"share_bg,omitempty"`
	RelationBg          string `json:"relation_bg,omitempty"`
	TaskBg              string `json:"task_bg,omitempty"`
	SupplierId          string `json:"supplier_id,omitempty"`            // 供应商ID
	ModelType           string `json:"model_type,omitempty"`             // 规格类型，单规格-single，多规格-multi
	IsCodeExchange      int8   `json:"is_code_exchange,omitempty"`       // 是否兑换码交易
	SubTitle            string `json:"sub_title,omitempty"`              // 子标题
	PageType            string `json:"page_type,omitempty"`              // newUser-新获客，gift-随手礼
	CanReceipt          int8   `json:"can_receipt,omitempty"`            // 是否支持开票
	MchId               string `json:"mch_id,omitempty"`                 // 微信收款商户号
	MemberTypePriceTip  string `json:"member_type_price_tip,omitempty"`  // 会员价说明
	OnshelfTime         string `json:"onshelf_time,omitempty"`           // 选了否，出现上架时间，自动上架
	IsActPoint          string `json:"is_act_point,omitempty"`           // 是否线上活动商品,0-无，1-线上，2-线下
	CouponStockId       string `json:"coupon_stock_id,omitempty"`        // 批次号
	CouponCustomerAppId string `json:"coupon_customer_app_id,omitempty"` // 应用号
	HasPointAct         int8   `json:"has_point_act,omitempty"`          // 是否有积分活动
	DigitalPayMchId     string `json:"digital_pay_mch_id,omitempty"`     // 数币收款商户号
	PayMethod           string `json:"pay_method,omitempty"`
}

func (this Goods) TableName() string {
	return "goods"
}

type UintArrayStruct []uint

func (AccessIds UintArrayStruct) Value() (driver.Value, error) {
	str := JsonMarshal(AccessIds)
	if str == "" {
		str = "[]"
	}
	return []byte(str), nil
}

func (AccessIds *UintArrayStruct) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	return json.Unmarshal(b, AccessIds)
}

type stringArrayStruct []string

func (AccessIds stringArrayStruct) Value() (driver.Value, error) {
	str := JsonMarshal(AccessIds)
	if str == "" {
		str = "[]"
	}
	return []byte(str), nil
}

func (AccessIds *stringArrayStruct) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	return json.Unmarshal(b, AccessIds)
}

func JsonMarshal(m interface{}) string {
	byteData, _ := json.Marshal(m)
	str := string(byteData)
	if str == "null" {
		str = ""
	}
	return str
}
