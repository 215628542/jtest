package model

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
)

var db *gorm.DB

func init() {
	var err error
	dsn := "root:@tcp(127.0.0.1:3305)/vmallnew?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err.Error())
	}

	// 数据表迁移
	if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").AutoMigrate(
		VmallInterPointCheckIn{}, VmallInterPointCheckInLog{}, VmallPointJob{}, VmallPointJobLog{}, VmallMemberPointLog{},
		VmallMember{}, VmallMemberPoint{}, Tt{},
	); err != nil {
		log.Println("表迁移失败：", err.Error())
	}

}

func GetDb() *gorm.DB {
	return db
}

type ActivityData struct {
	Id                  uint64    `json:"id" structs:"id" gorm:"primary_key;autoIncrement:false"`
	Cid                 uint64    `json:"cid" gorm:"index;not null;comment:活动id"`
	Uid                 string    `json:"uid" gorm:"index;not null;type:varchar(50);comment:用户id"`
	GoodsId             string    `json:"goods_id" gorm:"index;not null;type:varchar(50);comment:商品id"`
	WyfMerchantId       string    `json:"wyf_merchant_id" gorm:"index;not null;type:varchar(50);comment:商户号"`
	RuiyinxinMerchantId string    `json:"ruiyinxin_merchant_id" gorm:"index;not null;type:varchar(100);comment:瑞银信商户号"`
	RealnameDim         string    `json:"realname_dim" gorm:"not null;size:100;comment:真实姓名(带*号)"`
	RealnameEnc         string    `json:"realname_enc" gorm:"not null;size:200;comment:真实姓名(密文)"`
	CreateStaffId       string    `json:"create_staff_id" gorm:"not null;size:50;comment:创建员工"`
	MobileDim           string    `json:"mobile_dim" gorm:"not null;size:30;comment:手机号(带*号)"`
	MobileEnc           string    `json:"mobile_enc" gorm:"index;not null;size:100;comment:手机号(密文)"`
	IdCardDim           string    `json:"id_card_dim" gorm:"not null;size:100;comment:身份证(带*号)"`
	IdCardEnc           string    `json:"id_card_enc" gorm:"not null;size:100;comment:身份证(密文)"`
	OrgId               string    `json:"orgIds" gorm:"not null;size:50; comment:组织架构"`
	PrizeSpoil          string    `json:"prize_spoil" gorm:"not null;size:50;comment:奖品价值"`
	PrizeAmount         string    `json:"prize_amount" gorm:"not null;size:50;comment:奖品总价值"`
	PrizeName           string    `json:"prize_name" gorm:"not null;type:text;comment:奖品名称"`
	PrizeSum            string    `json:"prize_sum" gorm:"not null;size:50;comment:奖品数量"`
	State               bool      `json:"state" gorm:"index;not null;comment:领取状态"`
	TsState             string    `json:"ts_state" gorm:"index;not null;comment:推送状态"`
	OrderState          bool      `json:"order_state" gorm:"index;not null;comment:订单状态"`
	TimeDrawAt          time.Time `json:"time_draw_at" gorm:"not null;comment:领取时间"`
	Remark              string    `json:"remark" gorm:"not null;type:text;comment:备注"`
	BelongOrg           string    `json:"belong_org" gorm:"index,not null;size:50;comment:所属机构id"`
	OrderNum            string    `json:"order_num" gorm:"size:50;comment:id"`
	TaskId              string    `json:"task_id" gorm:"not null;size:100;comment:生成活动数据任务id"`

	ThirdReqBody string `json:"third_req_body" gorm:"not null;type:text;comment:请求报文"`
	ThirdRspBody string `json:"third_rsp_body" gorm:"not null;type:text;comment:响应报文"`
	ThirdChannel string `json:"third_channel" gorm:"not null;type:varchar(300);default:'';comment:三方通道"`

	ChannelId string    `json:"channel_id" gorm:"not null;index;type:varchar(50);comment:渠道ID"`
	DeletedAt time.Time `json:"deleted_at" gorm:"index;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null"`
}

func (ActivityData) TableName() string {
	return "activity_data"
}

type VmallMember struct {
	Id uint64 `json:"id" structs:"id" gorm:"primary_key;autoIncrement:false"`

	//Openid  string `json:"openid" structs:"openid" gorm:"not null;type:varchar(50)"`
	Appid   string `json:"appid" structs:"appid" gorm:"not null;index;type:varchar(50)"`
	Unionid string `json:"unionid" structs:"unionid"  gorm:"not null;type:varchar(50)"`
	Mobile  string `json:"mobile" structs:"mobile" gorm:"not null;index;type:varchar(30);comment:用户绑定的手机号"`

	Nickname string `json:"nickname" structs:"nickname" gorm:"not null;type:varchar(100);comment:用户昵称"`
	Avatar   string `json:"avatar" structs:"avatar"  gorm:"not null;type:varchar(300);commemt:头像"`
	Gendor   uint32 `json:"gendor" structs:"gendor" gorm:"not null;type:tinyint(2);comment:性别"`

	Level uint64 `json:"level" structs:"level" gorm:"not null;comment:vmall_levels表的level字段值"`

	ChannelId string `json:"channel_id" structs:"channel_id" gorm:"not null;type:varchar(50);comment:渠道ID;index:idx_channelid_deletedat_gdcstno,priority:1"`
	StaffId   string `json:"staff_id" structs:"staff_id" gorm:"not null;type:varchar(50);comment:归属员工id"`

	Realname string `json:"realname" structs:"realname" gorm:"not null;type:varchar(100);comment:真实姓名"`
	IdType   uint32 `json:"id_type"  gorm:"not null;type:tinyint(3);comment:证件种类:0-身份证 1-护照 2-港澳通行证 3-台湾同胞证"`
	Idcode   string `json:"idcode" gorm:"not null;type:varchar(50);comment:身份证"`
	Province string `json:"province" gorm:"not null;type:varchar(50);comment:省"`
	City     string `json:"city" gorm:"not null;type:varchar(100);comment:市"`
	District string `json:"district" gorm:"not null;type:varchar(100);comment:区"`
	Address  string `json:"address" gorm:"not null;type:varchar(300);comment:详细地址"`

	Longitude float64 `json:"longitude" gorm:"not null;comment:用户授权的经度"`
	Latitude  float64 `json:"latitude" gorm:"not null;comment:用户授权的维度"`

	OfflineMember     bool   `json:"offline_member" gorm:"not null;comment:是否线下会员-指代注册会员"`
	InviteStaffMobile string `json:"invite_staff_mobile" gorm:"not null;index;size:50;comment:推荐员工"`
	OrgId             string `json:"org_id" structs:"org_id" gorm:"not null;type:varchar(50);comment:员工关联的归属组织架构id"`

	IsPostalYsh bool `json:"is_postal_ysh" gorm:"not null;comment:是否邮生活会员"`

	KeepLocalUserType bool `json:"keep_local_user_type" gorm:"not null;comment:是否保持本地设置的会员身份"`

	SubId    string `json:"subId" gorm:"not null;size:50;comment:分站id"`
	SubOrgId string `json:"subOrgId" gorm:"not null;size:50;comment:分站网点id"` // 这个字段不用了，原因：注册时传入的值不准，所以都取subId做归属判断

	GzyzJoinErr string `json:"gzyz_join_err" gorm:"not null;type:text;comment:邮政入会异常信息记录"`

	GdOrgId   string `json:"gd_org_id" gorm:"not null;size:50;comment:广东省客管用户的网点id"`
	GdOrgNo   string `json:"gd_org_no" gorm:"not null;size:50;comment:广东省客管用户的网点机构号"`
	GdCstmNo  string `json:"gd_cstm_no" gorm:"not null;size:50;comment:广东省客管用户的客户机构号;index:idx_channelid_deletedat_gdcstno,priority:3"`
	GdCstmNos string `json:"gd_cstm_nos" gorm:"not null;size:200;comment:多个广东省客管用户的客户机构号"`

	DeletedAt time.Time `json:"deleted_at" structs:"deleted_at" gorm:"not null;index:idx_channelid_deletedat_gdcstno,priority:2"`
	CreatedAt time.Time `json:"created_at" structs:"created_at" gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at" structs:"updated_at" gorm:"not null"`
}

type VmallMemberPoint struct {
	Id uint64 `json:"id" gorm:"primary_key;autoIncrement:false"`

	Mobile    string `json:"mobile" gorm:"not null;type:varchar(50);index;comment:用户手机号;uniqueIndex:idx_mobile_channelid"`
	ChannelId string `json:"channel_id" gorm:"not null;type:varchar(50);index;comment:渠道id;uniqueIndex:idx_mobile_channelid"`

	Point uint64 `json:"point" gorm:"not null;comment:积分值"`

	CreatedAt time.Time `json:"created_at"  structs:"created_at"  gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at"  structs:"updated_at"  gorm:"not null"`
}

type VmallOrder struct {
	Id uint64 `json:"id" gorm:"primary_key;autoIncrement:false"`

	Uid           string `json:"uid" gorm:"not null;type:varchar(50);index;comment:member表的Id"`
	UserMobileDim string `json:"user_mobile_dim" gorm:"not null;size:50;comment:用户手机号(带*号)"`
	UserMobileEnc string `json:"user_mobile_enc" gorm:"index;not null;size:50;comment:用户手机号(密文)"`
	Appid         string `json:"appid" gorm:"not null;size:50;comment:appid"`
	Openid        string `json:"openid" gorm:"not null;size:50;comment:微信openid"`

	OrderNum      string `json:"order_num" gorm:"index;type:varchar(50);not null;comment:订单号"`
	TransactionId string `json:"transaction_id" gorm:"type:varchar(50);not null;comment:微信支付单号"`

	PayPoint        uint64 `json:"pay_point" gorm:"not null;comment:订单合计支付积分"`
	PayPrice        uint64 `json:"pay_price" gorm:"not null;comment:订单实际支付金额"`
	PayApiPointType string `json:"pay_api_point_type" gorm:"not null;type:varchar(50);comment:使用的接口积分类型"`
	PayApiPoint     uint64 `json:"pay_api_point" gorm:"not null;comment:使用接口积分"`

	State          string `json:"state" gorm:"not null;type:varchar(30);comment:订单状态;default:WAIT"`
	AfterSaleState string `json:"after_sale_state" gorm:"not null;type:varchar(30);comment:订单售后状态"`
	RefundRemark   string `json:"refund_remark" gorm:"not null;type:varchar(200);comment:记录订单异常状态"`

	Name        string `json:"name" gorm:"not null;type:varchar(50);comment:收件人姓名"`
	Mobile      string `json:"mobile" gorm:"not null;type:varchar(30);comment:收件人手机号"`
	Province    string `json:"province" gorm:"not null;type:varchar(50);comment:省"`
	City        string `json:"city" gorm:"not null;type:varchar(100);comment:市"`
	District    string `json:"district" gorm:"not null;type:varchar(100);comment:区"`
	Postcode    string `json:"postcode" gorm:"not null;type:varchar(20);comment:邮政编码"`
	Detail      string `json:"detail" gorm:"not null;type:varchar(300);comment:详细地址"`
	Logistics   string `json:"logistics" gorm:"not null;type:varchar(100);comment:物流号"`
	ExpressType string `json:"express_type" gorm:"not null;type:varchar(100);comment:快递类型"`

	ReverseOrgId  string    `json:"reverse_org_id" gorm:"not null;type:varchar(50);comment:预约网点id"`
	ReverseAt     time.Time `json:"reverse_at" gorm:"not null;comment:预约时间"`
	ReverseCode   string    `json:"reverse_code" gorm:"not null;type:varchar(6);comment:核销随机码"`
	WriteoffOrgId string    `json:"writeoff_org_id" gorm:"not null;type:varchar(50);comment:核销网点id"`
	WriteoffAt    time.Time `json:"writeoff_at" gorm:"not null;comment:核销时间"`

	IsLimitCard    uint32 `json:"is_limit_card" gorm:"not null;type:tinyint(1);comment:0不是限卡商品 1是限卡商品"`
	LimitCard      string `json:"limit_card" gorm:"not null;type:varchar(100);comment:限卡"`
	IsDistribution uint32 `json:"is_distribution" gorm:"not null;type:tinyint(1);comment:0不是分销订单 1是分销订单"`
	DistributorId  string `json:"distributor_id" gorm:"index;not null;type:varchar(50);comment:分销商会员id"`
	ThirdShopId    string `json:"third_shop_id" gorm:"not null;type:varchar(100);comment:三方的shop_id,有数据的时候核销需要推送信息"`
	ThirdOrderId   string `json:"third_order_id" gorm:"not null;type:varchar(100);comment:瑞银信返回的orderId"`
	ThirdReqBody   string `json:"third_req_body" gorm:"not null;type:text;comment:请求报文"`
	ThirdRspBody   string `json:"third_rsp_body" gorm:"not null;type:text;comment:响应报文"`
	ThirdSysErr    string `json:"third_sys_err" gorm:"not null;comment:调用异常"`

	PushGdPostalStatus uint32 `json:"push_gd_postal" gorm:"not null;comment:是否推送到省邮政接口 1-成功 2-失败"`
	PushGzPostalErr    string `json:"push_gz_req" gorm:"not null;type:text;comment:推送到省邮政接口错误信息"`

	ExchangeId    string `json:"exchange_id" gorm:"not null;type:varchar(50);comment:商品兑换id"`
	CanReceipt    bool   `json:"can_receipt" gorm:"not null;comment:是否能开票"`
	ReceiptStatus uint32 `json:"receipt_status" gorm:"not null;comment:开票状态 待开票-0 开票中-1 已开票-2"`

	ReceiveMchId string `json:"receive_mch_id" gorm:"not null;type:varchar(100);comment:收款商户id"`
	ReceiveMchNo string `json:"receive_mch_no" gorm:"not null;type:varchar(100);comment:收款商户号"`

	Remark string `json:"remark" gorm:"not null;size:400;comment:备注"`

	GiftPointActId string `json:"gift_point_act_id" gorm:"not null;size:50;comment:员工派积分活动id"`
	GiftPointOrgId string `json:"gift_point_org_id" gorm:"not null;size:50;comment:员工派发积分网点id"`
	GiftPointLogId string `json:"gift_point_log_id" gorm:"not null;size:50;comment:员工派积分记录Id"`
	ActivitySource string `json:"activity_source" gorm:"not null;size:50;comment:活动来源"`

	SysErr string `json:"sys_err" gorm:"not null;type:text;comment:错误信息"`

	ChannelId string    `json:"channel_id" structs:"channel_id" gorm:"not null;index;type:varchar(50);comment:渠道ID"`
	DelayId   string    `json:"delay_id" gorm:"not null;type:varchar(50);comment:"`
	DeletedAt time.Time `json:"deleted_at"   gorm:"not null"`
	CreatedAt time.Time `json:"created_at"   gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at"   gorm:"not null"`
	SuccessAt time.Time `json:"success_at" gorm:"not null;comment:订单完成时间"`
}

type GormStrSlice []string

func (tp GormStrSlice) Value() (driver.Value, error) {
	if len(tp) == 0 {
		return []byte("[]"), nil
	}

	return json.Marshal(tp)
}

func (tp *GormStrSlice) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	return json.Unmarshal(b, tp)
}

type GormJson []byte

func (tp GormJson) Value() (driver.Value, error) {
	if len(tp) == 0 {
		return []byte("{}"), nil
	}

	return tp, nil
}

func (tp *GormJson) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	*tp = b
	return nil
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

// 结构体转json字符串 抛弃异常
func JsonMarshal(m interface{}) string {
	byteData, _ := json.Marshal(m)
	str := string(byteData)
	if str == "null" {
		str = ""
	}
	return str
}

type VmallInterPointCheckIn struct {
	Id uint64 `json:"id,omitempty" gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL"`
	//SubIds      []string `json:"sub_ids,omitempty" gorm:"column:sub_ids;type:varchar(300);default:'';comment:'所属分站'"`
	CheckInType GormStrSlice `json:"check_in_type,omitempty" gorm:"column:check_in_type;type:text;comment:'签到类型'"`
	//CycleValue  []string  `json:"cycle_value,omitempty" gorm:"column:cycle_value;type:varchar(2000);default:'';comment:'周期对应积分数值'"`
	Cycle int32 `json:"cycle,omitempty" gorm:"column:cycle;type:int(11);default:0;comment:'签到周期'"`
	//Status     int32     `json:"status,omitempty" gorm:"column:status;type:int(11);default:0;comment:'状态  0：关闭 1:开启'"`
	RepCheckIn bool      `json:"rep_check_in,omitempty" gorm:"column:rep_check_in;type:int(11);default:0;comment:'中断签到是否重签'"`
	CreatedAt  time.Time `gorm:"column:created_at;default:NULL;comment:'创建时间'"`
	UpdatedAt  time.Time `gorm:"column:updated_at;default:NULL;comment:'更新时间'"`
	DeletedAt  time.Time `gorm:"column:deleted_at;default:NULL;comment:'删除时间'"`
}

type VmallInterPointCheckInLog struct {
	Id        int64     `json:"id" gorm:"primary_key;autoIncrement:false"`
	ActId     uint64    `json:"act_id,omitempty" gorm:"column:act_id;type:int(11);default:0;comment:'任务id'"`
	Uid       uint64    `json:"uid,omitempty" gorm:"column:uid;default:0;comment:'用户id'"`
	Point     int64     `json:"point,omitempty" gorm:"column:point;type:int(11);default:0;comment:'领取积分数'"`
	CreatedAt time.Time `gorm:"column:created_at;default:NULL;comment:'创建时间'"`
	UpdatedAt time.Time `gorm:"column:updated_at;default:NULL;comment:'更新时间'"`
	DeletedAt time.Time `gorm:"column:deleted_at;default:NULL;comment:'删除时间'"`
}

func (v *VmallInterPointCheckIn) Find(ctx context.Context, db *gorm.DB) []VmallInterPointCheckIn {

	l := []VmallInterPointCheckIn{}
	db.WithContext(ctx).Where("id=?", v.Id).Find(&l)
	return l
}

func (m *VmallPointJobLog) GetTodayReceiveNum(ctx context.Context, db *gorm.DB) (cnt int64, err error) {
	if m.Uid < 1 {
		err = errors.New("VmallPointJobLog 中的 uid 为空")
		return
	} else if m.Jid == 0 {
		err = errors.New("VmallPointJobLog 中的 Jid 为空")
		return
	}

	nowStart := fmt.Sprintf("%s 00:00:00", time.Now().Format("2006-01-02"))
	nowEnd := fmt.Sprintf("%s 23:59:59", time.Now().Format("2006-01-02"))
	err = db.WithContext(ctx).Model(VmallPointJobLog{}).Where("uid=? AND jid=? AND created_at BETWEEN ? AND ? AND status=?", m.Uid, m.Jid, nowStart, nowEnd, "RECEIVED").Count(&cnt).Error
	return
}

type VmallPointJob struct {
	Id uint64 `json:"id" gorm:"primary_key;autoIncrement:false"`

	NeedSub bool `json:"need_sub" gorm:"not null;comment:是否开启分站显示，false则全站显示"`

	Title   string `json:"title" gorm:"not null;type:varchar(100);comment:积分任务标题"`
	JobDesc string `json:"job_desc" gorm:"not null;type:text;comment:积分任务描述"'`
	IconImg string `json:"icon_img" gorm:"not null;type:varchar(50);comment:任务图标id"`
	NbLink  string `json:"nb_link" gorm:"not null;type:varchar(100);comment:内部链接"`

	Inactivated string `json:"inactivated" gorm:"not null;type:varchar(200);comment:未激活提示语"`
	IncrPoint   uint64 `json:"incr_point" gorm:"not null;comment:加分分值"`

	JobType string `json:"job_type" gorm:"not null;type:varchar(10);default:THIRD;comment: 三方-THIRD H5 公众号-OFFICAL"`

	GzPostalAct  uint64 `json:"gz_postal_act" gorm:"not null;comment:第三方专区系统活动ID"`
	GzPostalGift uint64 `json:"gz_postal_gift" gorm:"not null;comment:第三方专区系统礼品ID"`

	UrlLink string `json:"url_link" gorm:"not null;type:varchar(150);comment:H5的跳转链接"`

	OfficalAppid     string `json:"offical_appid" gorm:"not null;type:varchar(50);comment:公众号的appid"`
	OfficalAppsecret string `json:"offical_appsecret" gorm:"not null;type:varchar(50);comment:公众号的appsecret"`

	Status  uint32    `json:"status" gorm:"not null;comment:是否启用 1-是 0-否"`
	StartAt time.Time `json:"start_at" gorm:"not null;comment:开始时间"`
	EndAt   time.Time `json:"end_at" gorm:"not null;comment:结束时间"`

	ChannelId string `json:"channel_id" gorm:"not null;type:varchar(50);index;comment:渠道id"`

	// 任务类型-分享活动
	ShareNum      int32  `json:"share_num" gorm:"column:share_num;type:int(11);default:0;comment:每用户可完成次数"`
	ShareNumDaily int32  `json:"share_num_daily" gorm:"column:share_num_daily;type:int(11);default:0;comment:每用户每天可完成次数"`
	ActType       string `json:"act_type" gorm:"not null;type:varchar(50);default:'';comment: 活动类型"`

	DeletedAt time.Time `json:"deleted_at"  structs:"deleted_at"  gorm:"index;not null"`
	CreatedAt time.Time `json:"created_at"  structs:"created_at"  gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at"  structs:"updated_at"  gorm:"not null"`
}

type VmallPointJobLog struct {
	Id uint64 `json:"id" gorm:"primary_key;autoIncrement:false"`

	Uid uint64 `json:"uid" gorm:"index;not null;comment:member表的Id"`
	Jid uint64 `json:"jid" gorm:"index;not null;comment:job表的Id"`

	OrderNum string `json:"order_num" gorm:"index;type:varchar(50);not null;comment:订单号"`
	Idcode   string `json:"idcode" gorm:"not null;index;type:varchar(50);comment:身份证号"`

	Title     string `json:"title" gorm:"not null;type:varchar(100);comment:积分任务标题"`
	IncrPoint uint64 `json:"incr_point" gorm:"not null;comment:加分分值"`

	GzPostalAct  uint64 `json:"gz_postal_act" gorm:"not null;comment:第三方专区系统活动ID"`
	GzPostalGift uint64 `json:"gz_postal_gift" gorm:"not null;comment:第三方专区系统礼品ID"`
	CstmNo       string `json:"cstm_no" gorm:"not null;type:varchar(100);comment:客户号"`
	OrgId        string `json:"org_id" gorm:"not null;type:varchar(100);comment:9位机构号"`
	QueryReq     string `json:"query_req" gorm:"not null;type:text;comment:查询资格请求数据"`
	QueryResp    string `json:"query_resp" gorm:"not null;type:text;comment:查询资格响应数据"`
	PushStatus   string `json:"push_status" gorm:"not null;type:varchar(10);comment:通知邮政状态 未通知-WAIT 推送成功-SUCCESS 推送失败-FAIL"`
	PushReq      string `json:"push_req" gorm:"not null;type:text;comment:请求数据"`
	PushResp     string `json:"push_resp" gorm:"not null;type:text;comment:响应数据"`

	Status string `json:"status" gorm:"not null;type:varchar(10);comment:领取状态 待领取-WAIT 已领取-RECEIVED 领取失败-FAIL 领取异常-ABNORMAL"`
	Remark string `json:"remark" gorm:"not null;type:varchar(100);comment:备注"`

	ReceivedAt time.Time `json:"received_at" gorm:"not null;comment:领取时间"`
	CreatedAt  time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"not null"`
}

type VmallMemberPointLog struct {
	Id uint64 `json:"id" gorm:"primary_key;autoIncrement:false"`

	PointId   uint64 `json:"point_id" gorm:"not null;index;comment:points表的Id"`
	PointType string `json:"point_type" gorm:"not null;type:varchar(10);comment:积分类型：本地积分-LOCAL;index:idx_created_state_pointtype,priority:3"`
	InterType string `json:"inter_type" gorm:"not null;type:varchar(10);comment:接口积分类型：邮政-POSTAL"`

	EventDesc   string `json:"event_desc" gorm:"not null;type:varchar(300);comment:事件描述"`
	EventId     uint64 `json:"event_id" gorm:"not null;index;comment:events表的Id"`
	ChangePoint int64  `json:"change_point" gorm:"not null;comment:积分变化值"`
	State       string `json:"state" gorm:"not null;type:varchar(50);comment:审核状态 未审核-WAIT 审核失败-FAIL 审核成功-SUCCESS  撤销-REVOKE;index:idx_created_state_pointtype,priority:2"`

	IsThird      bool   `json:"is_third" gorm:"not null;comment:是否三方调用"`
	ChangeFrom   string `json:"change_from" gorm:"not null;size:70;comment:来源途径"`
	Appid        string `json:"appid" gorm:"index;not null;size:50;comment:系统内应用号"`
	ThirdOrderNo string `json:"third_order_no" gorm:"index;not null;size:50;comment:外部订单号"`
	RevokePoint  int64  `json:"revoke_point" gorm:"not null;comment:已撤销的积分值"`

	RyxShopId      string `json:"ryx_shop_id" gorm:"not null;size:50;comment:瑞银信商户号"`
	RyxReqBody     string `json:"ryx_req_body" gorm:"not null;comment:瑞银信请求数据"`
	RyxRspBody     string `json:"ryx_rsp_body" gorm:"not null;comment:瑞银信响应数据"`
	RyxOrderId     string `json:"ryx_order_id" gorm:"not null;size:50;comment:推送给瑞银信返回的orderId"`
	NotRecordLimit bool   `json:"not_record_limit" gorm:"not null;comment:不记录到限额 true-不记录"`

	DeletedAt time.Time `json:"deleted_at"  structs:"deleted_at"  gorm:"index;not null"`
	CreatedAt time.Time `json:"created_at"  structs:"created_at"  gorm:"not null;index:idx_created_state_pointtype,priority:1"`
	UpdatedAt time.Time `json:"updated_at"  structs:"updated_at"  gorm:"not null"`
}
