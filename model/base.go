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

	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(1)                  // 最大打开连接数
	sqlDB.SetMaxIdleConns(1)                  // 最大空闲连接数
	sqlDB.SetConnMaxLifetime(1 * time.Second) // 连接最大生命周期

	// 数据表迁移
	if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").AutoMigrate(
		EquityGoodsActPushRecord{}, RankSetting{}, ActivityBonusRank{}, SyhIssue{}, SyhActivityUser{}, SyhActivityReceive{}, SyhActivityTask{}, SyhActivity{}, WxMiniExtend{}, Activity{}, ActivityBonus{}, KycInfo{}, ApiLog{}, VmallMemberRelation{}, VmallMemberRelationExtand{}, VmallMemberRelationExtandLog{}, VmallMember{},
	); err != nil {
		log.Println("表迁移失败：", err.Error())
	}
}

type EquityGoodsActPushRecord struct {
	Id uint64 `json:"id" gorm:"primary_key;autoIncrement:false"`

	ReqData  string `json:"req_data" gorm:"not null;type:text;comment:请求信息"`
	RespData string `json:"resp_data" gorm:"not null;type:text;comment:响应信息"`
	ErrMsg   string `json:"err_msg" gorm:"not null;type:text;comment:异常信息"`
	Status   string `json:"status" gorm:"not null;type:varchar(50);default:'';comment:响应状态 success:成功 fail:失败"`
	Type     string `json:"type" gorm:"not null;size:50;comment:三方接口"`

	GoodsId       string `json:"goods_id" structs:"source" gorm:"not null;type:varchar(50);default:'';comment:商品id;"`
	EquityGoodsId string `json:"equity_goods_id" structs:"source" gorm:"not null;type:varchar(50);default:'';comment:云仓商品id;"`
	OrderNum      string `json:"order_num" structs:"source" gorm:"not null;type:varchar(50);default:'';comment:订单号;"`
	Source        string `json:"source" structs:"source" gorm:"not null;type:varchar(50);default:'';comment:来源;"`
	SourceId      string `json:"source_id" structs:"source_id" gorm:"not null;type:varchar(50);default:'';comment:来源id;"`

	CreatedAt time.Time `json:"createdAt" gorm:"index:idx_created_at;not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null"`
}

type ActivityBonusRank struct {
	Id uint64 `json:"id" structs:"id" gorm:"primary_key;autoIncrement:false"`

	MemberId        uint64 `json:"member_id" gorm:"not null;type:bigint(20);default:0;comment:会员id"`
	PromoActivityId uint64 `json:"promo_activity_id" gorm:"not null;type:bigint(20);default:0;comment:推广活动id"`
	PromoMemberId   uint64 `json:"promo_member_id" gorm:"not null;type:bigint(20);default:0;comment:推广会员ID"`
	PromoNum        int32  `json:"promo_num" gorm:"not null;type:int(11);default:0;comment:推广拉新人数"`
	Point           int32  `json:"point" gorm:"not null;type:int(11);default:0;comment:奖励积分"`
	Rank            int32  `json:"rank" gorm:"not null;type:int(11);default:0;comment:排行"`

	RankSettingId        uint64 `json:"rank_setting_id" gorm:"not null;type:bigint(20);default:0;comment:rank_setting表id"`
	PromoMemberMobileEnc string `json:"promo_member_mobile_enc" gorm:"not null;type:varchar(100);default:'';comment:推广会员手机"`
	//MerchantNum          string `json:"merchant_num" gorm:"not null;type:varchar(100);default:'';comment:商户号"`
	//MerchantName         string `json:"merchant_name" gorm:"not null;type:varchar(100);default:'';comment:商户名称"`
	//MerchantCategoryName string `json:"merchant_category_name" gorm:"not null;type:varchar(300);default:'';comment:商户分类名称"`

	//Bonus  int64  `json:"bonus" gorm:"not null;type:int(11);default:0;comment:奖励积分值"`
	Status string `json:"status" gorm:"not null;type:varchar(30);default:'';comment:派发奖励状态 待派发:pending 派发成功:success 派发失败:fail"`
	ErrMsg string `json:"err_msg" gorm:"not null;type:varchar(1000);default:'';comment:错误信息"`

	//ChannelId string `json:"channel_id" structs:"channel_id" gorm:"not null;type:varchar(50);default:'';comment:渠道ID;"`

	CreatedAt time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null"`
	DeletedAt time.Time `json:"deleted_at" gorm:"not null"`
}

type RankSetting struct {
	Id uint64 `json:"id" structs:"id" gorm:"primary_key;autoIncrement:false"`

	PromoActivityId uint64 `json:"promo_activity_id" gorm:"not null;type:bigint(20);default:0;comment:推广活动id"`
	Name            string `json:"name" gorm:"not null;type:varchar(200);default:'';comment:排行名称"`
	Rule            string `json:"rule" gorm:"type:text;comments:排行奖励规则"`
	PromoNum        int32  `json:"promo_num" gorm:"not null;type:int(11);default:0;comment:引流人数"`
	BonusCycle      string `json:"bonus_cycle" gorm:"not null;type:varchar(100);default:'';comment:排行奖励周期"`
	TopImage        string `json:"top_image" gorm:"not null;type:varchar(2000);default:'';comment:顶部图片"`

	EffectiveDate time.Time `json:"effective_date" gorm:"not null;comment:生效时间"`
	IsShow        bool      `json:"is_show" gorm:"not null;type:int(2);default:0;comment:是否展示"`
	Sort          int32     `json:"sort" gorm:"not null;type:bigint(20);default:0;comment:排序"`
	Status        bool      `json:"status" gorm:"not null;type:int(2);default:0;comment:状态"`

	StaffId string `json:"staff_id" gorm:"not null;type:varchar(100);default:'';comment:员工id"`

	CreatedAt time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null"`
	DeletedAt time.Time `json:"deleted_at" gorm:"not null"`
}

type ProjectDataPushRecord struct {
	Id        uint64 `json:"id" gorm:"primaryKey;autoIncrement:false"`
	ChannelId string `json:"channel_id" structs:"channel_id" gorm:"not null;type:varchar(50);default:'';comment:渠道ID;"`

	Params   string `json:"params" gorm:"not null;type:text;comment:请求参数"`
	ReqData  string `json:"req_data" gorm:"not null;type:text;comment:接口请求参数"`
	RespData string `json:"resp_data" gorm:"not null;type:text;comment:接口响应参数"`

	Type   string `json:"type" structs:"type" gorm:"not null;type:varchar(100);default:'';comment:类型;"`
	Status string `json:"status" structs:"status" gorm:"not null;type:varchar(30);default:'';comment:状态;"`
	ErrMsg string `json:"err_msg" gorm:"not null;type:text;comment:错误信息"`

	Source   string `json:"source" structs:"source" gorm:"not null;type:varchar(50);default:'';comment:来源;"`
	SourceId string `json:"source_id" structs:"source_id" gorm:"not null;type:varchar(50);default:'';comment:来源id;"`

	CreatedAt time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null"`
}

type ActivityBonus struct {
	Id uint64 `json:"id" structs:"id" gorm:"primary_key;autoIncrement:false"`

	PromoActivityId uint64 `json:"promo_activity_id" gorm:"not null;type:bigint(20);default:0;comment:推广活动id"`
	MemberType      int32  `json:"member_type" gorm:"not null;type:tinyint(2);default:0;comment:会员类型 1:新拓客户 2:存量客户"`
	PromoType       string `json:"promo_type" gorm:"not null;type:varchar(30);default:'';comment:推广类型 会员：member 商户：merchant"`
	PromoBonusType  string `json:"promo_bonus_type" gorm:"not null;type:varchar(30);default:'';comment:推广奖励类型 参与奖励:join 核销奖励:writeoff"`
	MemberId        uint64 `json:"member_id" gorm:"not null;type:bigint(20);default:0;comment:会员id"`

	JoinMemberId        uint64 `json:"join_member_id" gorm:"not null;type:bigint(20);default:0;comment:参与会员ID"`
	JoinMemberMobileEnc string `json:"join_member_mobile_enc" gorm:"not null;type:varchar(100);default:'';comment:参与会员手机"`
	OrgId               string `json:"org_id" gorm:"not null;type:varchar(100);default:'';comment:关联商户所属机构、会员所属机构，用promo_type区分"`

	PromoMemberId        uint64 `json:"promo_member_id" gorm:"not null;type:bigint(20);default:0;comment:推广会员ID"`
	PromoMemberMobileEnc string `json:"promo_member_mobile_enc" gorm:"not null;type:varchar(100);default:'';comment:推广会员手机"`
	RelationMerchantNum  string `json:"relation_merchant_num" gorm:"not null;type:varchar(100);default:'';comment:关联商户号"`
	RelationMerchantName string `json:"relation_merchant_name" gorm:"not null;type:varchar(100);default:'';comment:关联商户名称"`

	//RelationMerchantOrgId string `json:"relation_merchant_org_id" gorm:"not null;type:varchar(100);default:'';comment:关联商户所属机构"`

	Source       string `json:"source" gorm:"not null;type:varchar(100);default:'';comment:来源"`
	ActivityId   string `json:"activity_id" gorm:"not null;type:varchar(50);default:'';comment:活动id"`
	ActivityName string `json:"activity_name" gorm:"not null;type:varchar(100);default:'';comment:活动名称"`
	ReceiveId    string `json:"receive_id" gorm:"not null;type:varchar(50);default:'';comment:领取id"`
	ErrMsg       string `json:"err_msg" gorm:"not null;type:varchar(1000);default:'';comment:错误信息"`
	ErrCode      int32  `json:"err_code" gorm:"not null;type:int(11);default:0;comment:错误码"`

	Bonus  int64  `json:"bonus" gorm:"not null;type:int(11);default:0;comment:奖励积分值"`
	Status string `json:"status" gorm:"not null;type:varchar(30);default:'';comment:派发奖励状态 待派发:pending 派发成功:success 派发失败:fail"`

	BonusAt          time.Time `json:"bonus_at" gorm:"not null;comment:奖励时间;"`
	JoinActivityDate time.Time `json:"join_activity_date" gorm:"not null;comment:参与活动时间;"`
	WriteoffDate     time.Time `json:"writeoff_date" gorm:"not null;comment:核销时间;"`
	RegisterDate     time.Time `json:"register_date" gorm:"not null;comment:注册时间;"`
	ChannelId        string    `json:"channel_id" structs:"channel_id" gorm:"not null;type:varchar(50);default:'';comment:渠道ID;"`

	CreatedAt time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null"`
	DeletedAt time.Time `json:"deleted_at" gorm:"not null"`
}

type Operator struct {
	IsSuper       bool      `json:"isSuper" gorm:"not null;comment:操作员是否为超管"`
	AdminId       string    `json:"adminId" gorm:"not null;size:50;comment:操作员id"`
	MobileDim     string    `json:"mobileDim" gorm:"not null;size:30;comment:操作员手机号"`
	MobileEncrypt string    `json:"mobileEncrypt" gorm:"not null;size:30;comment:操作员手机号"`
	Name          string    `json:"name" gorm:"not null;size:200;comment:操作员名"`
	NameEncrypt   string    `json:"nameEncrypt" gorm:"not null;size:200;comment:操作员名"`
	AdminOrgId    string    `json:"adminOrgId" gorm:"not null;size:50;comment:操作员所属组织架构id"`
	OperatedAt    time.Time `json:"operatedAt" gorm:"not null;comment:操作时间"`
}

type SyhIssue struct {
	Id    uint64 `json:"id" gorm:"primary_key;auto_increment:false"`
	ActId uint64 `json:"act_id" gorm:"index;not null;comment:活动ID"`
	Operator

	TaskIds   GormStrSlice `json:"task_ids" gorm:"type:text;comment:派发任务id"`
	Uid       string       `json:"uid" gorm:"not null;size:50;comment:使用的用户id"`
	CreatedAt time.Time    `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time    `json:"updated_at" gorm:"not null"`
	DeletedAt time.Time    `json:"deleted_at" gorm:"not null"`
}

type SyhActivityUser struct {
	Id uint64 `json:"id" structs:"id" gorm:"primary_key;autoIncrement:false"`

	Aid        uint64       `json:"aid" gorm:"index;not null;index:idx_actId_uid,priority:1;comment:活动id"`
	AuthPrizes GormStrSlice `json:"auth_prizes" gorm:"not null;type:text;comment:授权的奖品任务"`

	Uid         string `json:"uid" gorm:"index;not null;size:50;index:idx_actId_uid,priority:2;comment:用户id"`
	RealnameDim string `json:"realname_dim" gorm:"not null;size:50;comment:用户姓名(带*号)"`
	RealnameEnc string `json:"realname_enc" gorm:"not null;size:70;comment:用户姓名(加密)"`
	MobileDim   string `json:"mobile_dim" gorm:"not null;size:50;comment:用户手机号(带*号)"`
	MobileEnc   string `json:"mobile_enc" gorm:"index;not null;size:70;comment:用户手机号(加密)"`
	IdcodeDim   string `json:"idcode_dim" gorm:"not null;size:70;comment:身份证(带*号)"`
	IdcodeEnc   string `json:"idcode_enc" gorm:"not null;size:100;comment:身份证(加密)"`
	Openid      string `json:"openid" gorm:"not null;size:50;comment:用户openid"`
	Appid       string `json:"appid" gorm:"not null;size:50;comment:领取的appid"`

	//MemberCreatedAt time.Time `json:"member_created_at" gorm:"not null;comment:会员入会时间"`
	//GenerateCodeAt time.Time `json:"generate_code_at" gorm:"not null;comment:生成码的时间"`

	StaffId        string `json:"staff_id" gorm:"index;not null;size:50;comment:员工id"`
	StaffOrg       string `json:"staff_org" gorm:"not null;size:50;comment:员工当时的网点"`
	StaffMobileDim string `json:"staff_mobile_dim" gorm:"not null;size:50;comment:员工手机号(带*号)"`
	StaffMobileEnc string `json:"staff_mobile_enc" gorm:"index;not null;size:70;comment:员工手机号(加密)"`

	ChannelId string `json:"channel_id" gorm:"index;not null;size:50;comment:渠道id"`

	DeletedAt time.Time `json:"deleted_at" structs:"deleted_at" gorm:"not null;index:idx_actId_uid,priority:3;"`
	CreatedAt time.Time `json:"created_at" structs:"created_at" gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at" structs:"updated_at" gorm:"not null"`
}

type SyhActivityReceive struct {
	Id uint64 `json:"id" structs:"id" gorm:"primary_key;autoIncrement:false"`

	MemberId    uint64 `json:"member_id" gorm:"not null;type:bigint(20);default:0;comment:会员id"`
	RealnameEnc string `json:"realname_enc" gorm:"not null;type:varchar(100);default:'';comment:客户姓名"`
	MobileEnc   string `json:"mobile_enc" gorm:"not null;type:varchar(100);default:'';comment:客户手机号码"`

	ActivityId    uint64 `json:"activity_id" gorm:"not null;type:bigint(20);default:0;comment:通关活动id"`
	TaskId        uint64 `json:"task_id" gorm:"not null;type:bigint(20);default:0;comment:任务id"`
	TaskType      string `json:"task_type" gorm:"not null;type:varchar(100);default:'';comment:任务类型"`
	BonusNum      int32  `json:"bonus_num" gorm:"not null;type:tinyint(2);default:0;comment:完成任务奖励抽奖次数"`
	ReceiveSource string `json:"receive_source" gorm:"not null;type:varchar(100);default:'';comment:领取来源"`

	StaffId        string `json:"staff_id" gorm:"not null;type:varchar(100);default:'';comment:员工id"`
	StaffNameEnc   string `json:"staff_name_enc" gorm:"not null;type:varchar(100);default:'';comment:员工名称"`
	StaffOrgId     string `json:"staff_org_id" gorm:"not null;type:varchar(100);default:'';comment:员工机构id"`
	StaffMobileEnc string `json:"staff_mobile_enc" gorm:"not null;type:varchar(100);default:'';comment:员工手机号码"`

	ReceiveStatus string    `json:"receive_status" gorm:"not null;type:varchar(100);default:'';comment:领取状态"`
	ReceiveAt     time.Time `json:"receive_at" gorm:"not null;comment:领取时间"`
	VerifyDate    time.Time `json:"verify_date" gorm:"not null;comment:验证时间"`

	CreatedAt time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null"`
	DeletedAt time.Time `json:"deleted_at" gorm:"not null"`
}

type SyhActivityTask struct {
	Id uint64 `json:"id" structs:"id" gorm:"primary_key;autoIncrement:false"`

	Name          string `json:"name" gorm:"not null;type:varchar(300);default:'';comment:任务名称"`
	ActivityId    uint64 `json:"activity_id" gorm:"not null;type:bigint(20);default:0;comment:通关活动id"`
	TaskType      string `json:"task_type" gorm:"not null;type:varchar(300);default:'';comment:任务类型"`
	BonusNum      int32  `json:"bonus_num" gorm:"not null;type:tinyint(2);default:0;comment:完成任务奖励抽奖次数"`
	IsShare       bool   `json:"is_share" gorm:"not null;type:int(2);default:0;comment:是否允许分享"`
	IsStaffUnlock bool   `json:"is_staff_unlock" gorm:"not null;type:int(2);default:0;comment:启用员工端解锁"`

	UnFinishTips    string `json:"un_finish_tips" gorm:"not null;type:varchar(500);default:'';comment:未办理业务提示语"`
	ExplainLink     string `json:"explain_link" gorm:"not null;type:varchar(500);default:'';comment:业务介绍链接"`
	UnFinishImageId string `json:"un_finish_image_id" gorm:"not null;type:varchar(500);default:'';comment:任务已完成图片"`
	FinishImageId   string `json:"finish_image_id" gorm:"not null;type:varchar(500);default:'';comment:任务未完成图片"`

	Status bool  `json:"status" gorm:"not null;type:int(2);default:1;comment:状态"`
	Sort   int32 `json:"sort" gorm:"not null;type:int(11);default:0;comment:排序"`

	CreatedAt time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null"`
	DeletedAt time.Time `json:"deleted_at" gorm:"not null"`
}

type SyhActivity struct {
	Id uint64 `json:"id" structs:"id" gorm:"primary_key;autoIncrement:false"`

	Name               string    `json:"name" gorm:"not null;type:varchar(300);default:'';comment:活动名称"`
	OrgId              string    `json:"org_id" gorm:"not null;type:varchar(100);default:'';comment:所属机构id"`
	StartDate          time.Time `json:"start_date" gorm:"not null;comment:活动开始时间"`
	EndDate            time.Time `json:"end_date" gorm:"not null;comment:活动结束时间"`
	RelationActivityId string    `json:"relation_activity_id" gorm:"not null;type:varchar(500);default:'';comment:关联活动ID"`
	Rule               string    `json:"rule" gorm:"type:text;comments:活动规则"`
	VerifyDate         time.Time `json:"start_date" gorm:"not null;comment:验证时间"`

	IsShare      bool   `json:"is_share" gorm:"not null;type:int(2);default:0;comment:是否允许分享"`
	ShareImageId string `json:"share_image_id" gorm:"not null;type:varchar(300);default:'';comment:分享图"`
	ShareTitle   string `json:"share_title" gorm:"not null;type:varchar(500);default:'';comment:分享标题"`

	BgImageId      string `json:"bg_image_id" gorm:"not null;type:varchar(300);default:'';comment:背景图片"`
	GridImageId    string `json:"grid_image_id" gorm:"not null;type:varchar(300);default:'';comment:宫格图"`
	PosterImageId  string `json:"poster_image_id" gorm:"not null;type:varchar(300);default:'';comment:海报图"`
	NotTaskImageId string `json:"not_task_image_id" gorm:"not null;type:varchar(300);default:'';comment:无任务格图"`

	Status    bool   `json:"status" gorm:"not null;type:int(2);default:1;comment:状态"`
	ChannelId string `json:"channel_id" structs:"channel_id" gorm:"not null;type:varchar(50);default:'';comment:渠道ID;"`

	CreatedAt time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null"`
	DeletedAt time.Time `json:"deleted_at" gorm:"not null"`
}

type VmallInterPointCheckInLog struct {
	Id        int64  `json:"id" gorm:"primary_key;autoIncrement:false"`
	ActId     uint64 `json:"act_id,omitempty" gorm:"column:act_id;type:int(11);default:0;comment:'任务id'"`
	Uid       uint64 `json:"uid,omitempty" gorm:"column:uid;default:0;comment:'用户id'"`
	Point     int64  `json:"point,omitempty" gorm:"column:point;type:int(11);default:0;comment:'领取积分数'"`
	PointType string `json:"point_type" gorm:"column:point_type;type:varchar(50);default:'';comment:'积分类型'"`

	CreatedAt time.Time `gorm:"column:created_at;default:NULL;comment:'创建时间'"`
	UpdatedAt time.Time `gorm:"column:updated_at;default:NULL;comment:'更新时间'"`
	DeletedAt time.Time `gorm:"column:deleted_at;default:NULL;comment:'删除时间'"`
}

type VmallInterPointCheckIn struct {
	Id uint64 `json:"id" gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL"`
	//SubIds      db.StringArrayStruct `json:"sub_ids" gorm:"column:sub_ids;type:text;comment:'所属分站'"`
	CheckInType string `json:"check_in_type" gorm:"column:check_in_type;type:varchar(100);default:'';comment:'签到类型'"`
	//CycleValue  db.StringArrayStruct `json:"cycle_value" gorm:"column:cycle_value;type:text;default:'';comment:'周期对应积分数值'"`
	Cycle      int32 `json:"cycle" gorm:"column:cycle;type:int(11);default:0;comment:'签到周期'"`
	Status     int32 `json:"status" gorm:"column:status;type:int(11);default:0;comment:'状态  0：关闭 1:开启'"`
	RepCheckIn bool  `json:"rep_check_in" gorm:"column:rep_check_in;type:int(11);default:0;comment:'中断签到是否重签'"`

	PointType string `json:"point_type" gorm:"column:point_type;type:varchar(50);default:'';comment:'积分类型'"`
	SubType   string `json:"sub_type" gorm:"column:sub_type;type:varchar(50);default:'';comment:'分站应用范围'"`

	CreatedAt time.Time `gorm:"column:created_at;default:NULL;comment:'创建时间'"`
	UpdatedAt time.Time `gorm:"column:updated_at;default:NULL;comment:'更新时间'"`
	DeletedAt time.Time `gorm:"column:deleted_at;default:NULL;comment:'删除时间'"`
}

type WxMiniExtend struct {
	gorm.Model

	ID                    uint   `gorm:"primarykey;autoIncrement:false"`
	MiniID                uint   `gorm:"column:mini_id;index:index_wme_miniid;not null;comment:wx_minis表的主键ID"`
	ChannelID             string `gorm:"column:channel_id;type:varchar(50);comment:渠道ID"`
	IssueMchID            string `gorm:"column:issue_mchid;type:varchar(50);comment:发券方商户ID"`
	PayMchID              string `gorm:"column:pay_mchid;type:varchar(50);comment:企业付款商户ID"`
	RecMchID              string `gorm:"column:rec_mchid;type:varchar(50);comment:企业收款商品ID"`
	BindcardMchID         string `gorm:"column:bindcard_mchid;type:varchar(50);comment:绑卡验卡商户号"`
	RegisterMemberType    uint32 `json:"register_member" gorm:"not null;type:tinyint(1);comment:是否开启会员注册 0关闭 1开启"`
	MemberTypeCode        string `json:"member_type_id" gorm:"not null;type:varchar(50);comment:会员注册对应的会员类型编码"`
	RegisterSuccessTip    string `json:"register_success_tip" gorm:"not null;type:varchar(200);comment:注册成功提示"`
	OfficalAccountLink    string `json:"offical_account_link" gorm:"not null;type:varchar(200);comment:公众号链接"`
	UserServicesAgreement string `json:"user_services_agreement" gorm:"not null;type:text;comment:用户服务协议"`
	UserPrivacyPolicy     string `json:"user_privacy_policy" gorm:"not null;type:text;comment:用户隐私条例"`
	ActReceiveMchId       string `json:"act_receive_mch_id" gorm:"not null;type:varchar(50);comment:活动收款商户号"`
	DigitalMchNo          string `json:"digital_mch_no" gorm:"not null;type:varchar(50);comment:数币商户号"`

	RegisterInfo  GormStrSlice `json:"register_info" gorm:"not null;type:json;comment:小程序注册必填字段"`
	TitleBarColor string       `json:"title_bar_color" gorm:"not null;type:varchar(50);comment:标题栏颜色"`
	TitleBarText  string       `json:"title_bar_text" gorm:"not null;type:varchar(100);comment:标题栏标题文本"`
	BgImgId       string       `json:"bg_img_id" gorm:"not null;type:varchar(100);comment:背景图"`
	NeedPushThird bool         `json:"need_push_third" gorm:"not null;comment:是否推送三方接口"`
	BindInfoTip   string       `json:"bind_info_tip" gorm:"not null;type:varchar(200);comment:绑定信息提示语"`
	LoginTip      string       `json:"login_tip" gorm:"not null;type:varchar(200);comment:登录提示语"`

	OrderAutoReceive    bool   `json:"order_auto_receive" gorm:"not null;comment:是否开启自动收货"`
	OrderAutoReceiveDay uint32 `json:"order_auto_receive_day" gorm:"not null;comment:自动收货的天数"`

	ExtUrl string `json:"ext_url" gorm:"not null;size:200;comment:客服链接"`
	CorpId string `json:"corp_id" gorm:"not null;size:50;comment:企业id"`

	IndexBannerImgId string `json:"index_banner_img_id" gorm:"not null;size:50;comment:首页banner图id"`

	SupportOnePointBuy bool   `json:"support_one_point_buy" gorm:"not null;comment:员工端是否支持创建一分购类型代销商品"`
	OfficialAppId      string `json:"official_app_id" gorm:"not null;size:50;comment:公众号appid"`
	OfficialAppSecret  string `json:"official_app_secret" gorm:"not null;size:100;comment:公众号appsecret"`

	AliMchNo                  string `json:"ali_mch_no" gorm:"column:ali_mch_no;not null;type:varchar(100);default:'';comment:支付宝商户号"`
	RedBagSendUserDayLimitWx  uint64 `json:"red_bag_send_user_day_limit_wx" gorm:"column:red_bag_send_user_day_limit_wx;not null;type:varchar(100);default:0;comment:红包单用户每天派发限额-微信"`
	RedBagSendUserDayLimitAli uint64 `json:"red_bag_send_user_day_limit_ali" gorm:"column:red_bag_send_user_day_limit_ali;not null;type:varchar(100);default:0;comment:红包单用户每天派发限额-支付宝"`

	SupplierAdmissionTerms string `json:"supplier_admission_terms" gorm:"type:text;comment:供应商入驻条款"`

	HomeIrregularImgId   string `json:"home_irregular_img_id" gorm:"not null;type:varchar(100);default:'';comment:用户端小程序首页异形图片id"`
	HomeIrregularBgColor string `json:"home_irregular_bg_color" gorm:"not null;type:varchar(100);default:'';comment:用户端小程序首页异形图片背景颜色"`

	IsOpenStaffRegisterCode bool         `json:"is_open_staff_register_code" gorm:"column:is_open_staff_register_code;not null;type:tinyint(1);default:0;comment:是否开启员工端注册码入会附表"`
	UserkOrgIds             GormStrSlice `json:"userk_org_ids" gorm:"not null;type:json;comment:首K适用机构"`
	MemberTags              GormStrSlice `json:"member_tags" gorm:"not null;type:json;comment:附表已完成用户标签"`
	UserkEquityBagId        string       `json:"userk_equity_bag_id" gorm:"column:userk_equity_bag_id;not null;type:varchar(300);default:'';comment:首K关联权益礼包ID"`
	UserkJumpType           string       `json:"userk_jump_type" gorm:"column:userk_jump_type;not null;type:varchar(500);default:'';comment:商户KYC受用机构"`
	UserkRecipientOrgIds    GormStrSlice `json:"userk_recipient_org_ids" gorm:"not null;type:json;comment:商户KYC填写后跳转链接类型"`
}

type Activity struct {
	Id uint64 `json:"id" structs:"id" gorm:"primary_key;autoIncrement:false"`

	Name               string    `json:"name" gorm:"not null;type:varchar(300);default:'';comment:活动名称"`
	OrgId              string    `json:"org_id" gorm:"not null;type:varchar(100);default:'';comment:所属机构"`
	StartDate          time.Time `json:"start_date" gorm:"not null;comment:活动开始时间"`
	EndDate            time.Time `json:"end_date" gorm:"not null;comment:活动结束时间"`
	AdImage            string    `json:"ad_image" gorm:"type:text;comments:广告图"`
	ActivityType       int32     `json:"activity_type" gorm:"not null;type:tinyint(2);default:0;comment:活动类型 1:注册类 2:活动类"`
	Sort               int32     `json:"sort" gorm:"not null;type:int(11);default:0;comment:排序"`
	DistanceConstraint int32     `json:"distance_constraint" gorm:"not null;type:int(11);default:0;comment:距离约束（米）"`

	RegisterBonus             int64 `json:"register_bonus" gorm:"not null;type:int(11);default:0;comment:新拓客户注册奖励积分（选注册类时出现该选项必选）"`
	EligibleUserRegisterBonus int64 `json:"eligible_user_register_bonus" gorm:"not null;type:int(11);default:0;comment:存量客户注册奖励积分（选注册类时出现该选项必选）"`
	NewRegister30DayBonus     int64 `json:"new_register_30_day_bonus" gorm:"not null;type:int(11);default:0;comment:30天内新开户奖励积分（选注册类或活动类时都出现该选项非必选）"`

	NotEligibleUserJoinBonus int64 `json:"not_eligible_user_join_bonus" gorm:"not null;type:int(11);default:0;comment:新拓客户参与奖励积分（选活动类时出现该选项必选）"`
	EligibleUserJoinBonus    int64 `json:"eligible_user_join_bonus" gorm:"not null;type:int(11);default:0;comment:存量客户参与奖励积分（选活动类时出现该选项必选）"`

	NotEligibleUserWriteoffBonus int64 `json:"not_eligible_user_writeoff_bonus" gorm:"not null;type:int(11);default:0;comment:新拓客户核销奖励积分（选活动类时出现该选项非必选）"`

	ActivityLink     string `json:"activity_link" gorm:"not null;type:varchar(300);default:'';comment:活动链接"`
	PromoCrowds      string `json:"promo_crowds" gorm:"not null;type:varchar(300);default:'';comments:可推广人群 商户:merchant 会员:member"`
	RelationMerchant int32  `json:"relation_merchant" gorm:"not null;type:tinyint(2);default:1;comment:关联商户 1:全量推广商户 2:单独关联商户"`
	Status           int32  `json:"status" gorm:"not null;type:tinyint(2);default:1;comment:1:开启 2:关闭"`

	StaffId    string `json:"staff_id" gorm:"not null;type:varchar(100);default:'';comment:员工id"`
	StaffOrgId string `json:"staff_org_id" gorm:"not null;type:varchar(100);default:'';comment:员工组织架构id"`
	ChannelId  string `json:"channel_id" structs:"channel_id" gorm:"not null;type:varchar(50);default:'';comment:渠道ID;"`

	CreatedAt time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null"`
	DeletedAt time.Time `json:"deleted_at" gorm:"not null"`
}

type KycInfo struct {
	Id uint64 `json:"id" structs:"id" gorm:"primary_key;autoIncrement:false"`

	StaffId    string `json:"staff_id" gorm:"not null;type:varchar(100);default:'';comment:员工id"`
	StaffOrgId string `json:"staff_org_id" gorm:"not null;type:varchar(100);default:'';comment:员工组织架构id"`
	ChannelId  string `json:"channel_id" structs:"channel_id" gorm:"not null;type:varchar(50);default:'';comment:渠道ID;"`

	ProjectName             string `json:"project_name" gorm:"not null;type:varchar(100);default:'';comment:项目名称"`
	ProjectManagerMobileEnc string `json:"project_manager_mobile_enc" gorm:"not null;type:varchar(100);default:'';comment:项目负责人姓名手机号(加密)"`
	ProjectManagerNameEnc   string `json:"project_manager_name_enc" gorm:"not null;type:varchar(100);default:'';comment:项目负责人姓名(加密)"`
	ManagerOrgId            string `json:"manager_org_id" gorm:"not null;type:varchar(300);default:'';comment:项目负责人组织架构"`
	DetailAddress           string `json:"detail_address" gorm:"not null;type:varchar(300);default:'';comment:详细地址"`
	ProjectAddress          string `json:"project_address" gorm:"not null;type:varchar(300);default:'';comment:项目地址"`
	DetailAddressLongitude  string `json:"detail_address_longitude" gorm:"not null;type:varchar(100);default:'';comment:详细地址经度"`
	DetailAddressLatitude   string `json:"detail_address_latitude" gorm:"not null;type:varchar(100);default:'';comment:详细地址纬度"`

	ManagementUnit     string `json:"management_unit" gorm:"not null;type:varchar(100);default:'';comment:管理方单位"`
	KycType            string `json:"kyc_type" gorm:"not null;type:varchar(100);default:'';comment:行业类型"`
	TransactionSize    string `json:"transaction_size" gorm:"not null;type:varchar(100);default:'';comment:交易规模"`
	MerchantNum        int32  `json:"merchant_num" gorm:"not null;type:int(11);default:0;comment:商户数"`
	IsHasIndustryChain bool   `json:"is_has_industry_chain" gorm:"not null;type:int(2);default:0;comment:是否有产业链"`
	UpMerchantType     string `json:"transaction_size" gorm:"not null;type:varchar(100);default:'';comment:上游商户类型"`
	UpMerchantOrg      string `json:"transaction_size" gorm:"not null;type:varchar(100);default:'';comment:上游商户分布区域"`
	DownMerchantType   string `json:"transaction_size" gorm:"not null;type:varchar(100);default:'';comment:下游商户类型"`
	DownMerchantOrg    string `json:"transaction_size" gorm:"not null;type:varchar(100);default:'';comment:下游商户分布区域"`

	PaymentTool     string `json:"payment_tool" gorm:"not null;type:varchar(100);default:'';comment:支付工具"`
	CooperativeBank string `json:"cooperative_bank" gorm:"not null;type:varchar(100);default:'';comment:合作银行"`
	FinancialNeeds  string `json:"financial_needs" gorm:"not null;type:varchar(100);default:'';comment:金融需求"`
	FinancialMs     string `json:"financial_ms" gorm:"not null;type:varchar(100);default:'';comment:主要财务管理系统"`

	OrgId              stringArrayStruct `json:"images" gorm:"type:text; comments:归口网点"`
	TeamManagerNameEnc string            `json:"team_manager_name_enc" gorm:"not null;type:varchar(100);default:'';comment:开发团队负责人"`
	TeamMembers        string            `json:"team_members" gorm:"not null;type:varchar(300);default:'';comment:开发团队成员"`

	// 项目跟进信息
	ClockInLongitude    string            `json:"clock_in_longitude" gorm:"not null;type:varchar(100);default:'';comment:打卡经度"`
	ClockInLatitude     string            `json:"clock_in_latitude" gorm:"not null;type:varchar(300);default:'';comment:打卡纬度"`
	ClockInAddress      string            `json:"clock_in_address" gorm:"not null;type:varchar(300);default:'';comment:打卡位置"`
	IsOpenClockIn       bool              `json:"is_open_clock_in" gorm:"not null;type:int(2);default:0;comment:是否开启打卡位置"`
	ProjectFollowDesc   string            `json:"project_follow_desc" gorm:"not null;type:varchar(1000);default:'';comment:项目跟进描述"`
	ProjectFollowDate   time.Time         `json:"created_at" gorm:"not null;comment:项目跟进日期"`
	ProjectFollowImages stringArrayStruct `json:"project_follow_images" gorm:"type:text;comments:目跟进图片"`

	CreatedAt time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null"`
	DeletedAt time.Time `json:"deleted_at" gorm:"not null"`
}

type ApiLog struct {
	Id             uint64 `json:"id" gorm:"primary_key;autoIncrement:false"`
	Param          string `json:"param" gorm:"not null;type:text;comment:服务请求参数"`
	ReqData        string `json:"req_data" gorm:"not null;type:text;comment:请求信息"`
	RespData       string `json:"resp_data" gorm:"not null;type:text;comment:响应信息"`
	OriginRespData string `json:"origin_resp_data" gorm:"not null;type:text;comment:原始响应信息"`

	ErrMsg     string `json:"err_msg" gorm:"not null;type:text;comment:异常信息"`
	StatusCode string `json:"status_code" gorm:"not null;default:'';comment:响应状态码 00为验证成功，01客户信息验证失败，99为传入参数不合法"`

	Unionid   string `json:"unionid" gorm:"not null;size:100;default:'';comment:unionid"`
	MobileEnc string `json:"mobile_enc" gorm:"not null;size:150;default:'';comment:加密手机号码"`

	Source    string `json:"source" gorm:"not null;index:idx_source;size:100;default:'';comment:来源"`
	ServiceId string `json:"service_id" gorm:"not null;size:50;default:'';comment:接口服务id"`
	ChannelId string `json:"channel_id" gorm:"not null;size:50;default:'';comment:渠道ID"`

	CreatedAt time.Time `gorm:"column:created_at;default:NULL;comment:'创建时间'"`
	UpdatedAt time.Time `gorm:"column:updated_at;default:NULL;comment:'更新时间'"`
}

type VmallMember struct {
	Id uint64 `json:"id" structs:"id" gorm:"primary_key;autoIncrement:false"`

	//Openid  string `json:"openid" structs:"openid" gorm:"not null;type:varchar(50)"`
	Appid   string `json:"appid" structs:"appid" gorm:"not null;index;type:varchar(50)"`
	Unionid string `json:"unionid" structs:"unionid"  gorm:"not null;type:varchar(50)"`

	//Mobile  string `json:"mobile" structs:"mobile" gorm:"not null;index;type:varchar(30);comment:用户绑定的手机号"`
	//Realname string `json:"realname" structs:"realname" gorm:"not null;type:varchar(100);comment:真实姓名"`
	//Idcode   string `json:"idcode" gorm:"not null;type:varchar(50);comment:身份证"`
	//InviteStaffMobile string `json:"invite_staff_mobile" gorm:"not null;index;size:50;comment:推荐员工"`
	MobileDim                string `json:"mobile_dim" structs:"mobile_dim" gorm:"not null;index;type:varchar(30);comment:用户绑定的手机号"`
	MobileEncrypt            string `json:"mobile_encrypt" structs:"mobile_encrypt" gorm:"not null;index;type:varchar(30);comment:用户绑定的手机号"`
	RealnameDim              string `json:"realname_dim" structs:"realname_dim" gorm:"not null;type:varchar(100);comment:真实姓名"`
	RealnameEncrypt          string `json:"realname_encrypt" structs:"realname_encrypt" gorm:"not null;type:varchar(100);comment:真实姓名"`
	IdCodeDim                string `json:"id_code_dim" gorm:"not null;type:varchar(50);comment:身份证"`
	IdCodeEncrypt            string `json:"id_code_encrypt" gorm:"not null;type:varchar(50);comment:身份证"`
	InviteStaffMobileDim     string `json:"invite_staff_mobile_dim" gorm:"not null;index;size:50;comment:推荐员工"`
	InviteStaffMobileEncrypt string `json:"invite_staff_mobile_encrypt" gorm:"not null;index;size:50;comment:推荐员工"`

	Nickname string `json:"nickname" structs:"nickname" gorm:"not null;type:varchar(100);comment:用户昵称"`
	Avatar   string `json:"avatar" structs:"avatar"  gorm:"not null;type:varchar(300);commemt:头像"`
	Gendor   uint32 `json:"gendor" structs:"gendor" gorm:"not null;type:tinyint(2);comment:性别"`

	Level            uint64 `json:"level" structs:"level" gorm:"not null;comment:vmall_levels表的level字段值"`
	IsSyhOldCustomer bool   `json:"is_syh_new_customer" gorm:"not null;type:int(2);default:0;comment:是否苏邮惠老客户"`

	ChannelId string `json:"channel_id" structs:"channel_id" gorm:"not null;type:varchar(50);comment:渠道ID;index:idx_channelid_deletedat_gdcstno,priority:1"`
	StaffId   string `json:"staff_id" structs:"staff_id" gorm:"not null;type:varchar(50);comment:归属员工id"`

	IdType   uint32 `json:"id_type"  gorm:"not null;type:tinyint(3);comment:证件种类:0-身份证 1-护照 2-港澳通行证 3-台湾同胞证"`
	Province string `json:"province" gorm:"not null;type:varchar(50);comment:省"`
	City     string `json:"city" gorm:"not null;type:varchar(100);comment:市"`
	District string `json:"district" gorm:"not null;type:varchar(100);comment:区"`
	Address  string `json:"address" gorm:"not null;type:varchar(300);comment:详细地址"`

	Longitude float64 `json:"longitude" gorm:"not null;comment:用户授权的经度"`
	Latitude  float64 `json:"latitude" gorm:"not null;comment:用户授权的维度"`

	OfflineMember bool   `json:"offline_member" gorm:"not null;comment:是否线下会员-指代注册会员"`
	OrgId         string `json:"org_id" structs:"org_id" gorm:"not null;type:varchar(50);comment:员工关联的归属组织架构id"`

	IsPostalYsh bool `json:"is_postal_ysh" gorm:"not null;comment:是否邮生活会员"`

	KeepLocalUserType bool `json:"keep_local_user_type" gorm:"not null;comment:是否保持本地设置的会员身份"`

	SubId    string `json:"subId" gorm:"not null;size:50;comment:分站id"`
	SubOrgId string `json:"subOrgId" gorm:"not null;size:50;comment:分站网点id"` // 这个字段不用了，原因：注册时传入的值不准，所以都取subId做归属判断

	GzyzJoinErr string `json:"gzyz_join_err" gorm:"not null;type:text;comment:邮政入会异常信息记录"`

	GdOrgId   string `json:"gd_org_id" gorm:"not null;size:50;comment:广东省客管用户的网点id"`
	GdOrgNo   string `json:"gd_org_no" gorm:"not null;size:50;comment:广东省客管用户的网点机构号"`
	GdCstmNo  string `json:"gd_cstm_no" gorm:"not null;size:50;comment:广东省客管用户的客户机构号;index:idx_channelid_deletedat_gdcstno,priority:3"`
	GdCstmNos string `json:"gd_cstm_nos" gorm:"not null;size:200;comment:多个广东省客管用户的客户机构号"`

	YshUserId string `json:"yshUserId" gorm:"not null;size:50;comment:邮生活的用户id"`

	DeletedAt time.Time `json:"deleted_at" structs:"deleted_at" gorm:"not null;index:idx_channelid_deletedat_gdcstno,priority:2"`
	CreatedAt time.Time `json:"created_at" structs:"created_at" gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at" structs:"updated_at" gorm:"not null"`
}

type VmallMemberRelation struct {
	Id uint64 `json:"id" structs:"id" gorm:"primary_key;autoIncrement:false"`

	CustomerNameEnc   string `json:"customer_name" gorm:"not null;type:varchar(100);default:'';comment:客户名(加密)"`
	CustomerMobileEnc string `json:"customer_mobile_enc" gorm:"not null;type:varchar(100);default:'';comment:客户手机号(加密)"`
	CustomerIdCodeEnc string `json:"customer_id_code_enc" gorm:"not null;type:varchar(100);default:'';comment:客户身份证(加密)"`

	StaffId    string `json:"staff_id" gorm:"not null;type:varchar(100);default:'';comment:员工id"`
	StaffOrgId string `json:"staff_org_id" gorm:"not null;type:varchar(100);default:'';comment:员工组织架构id"`
	ChannelId  string `json:"channel_id" structs:"channel_id" gorm:"not null;type:varchar(50);default:'';comment:渠道ID;"`

	CreatedAt time.Time `json:"createdAt" gorm:"not null"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"not null"`
	DeletedAt time.Time `json:"deleted_at" gorm:"not null"`
}

type VmallMemberRelationExtand struct {
	Id uint64 `json:"id" structs:"id" gorm:"primary_key;autoIncrement:false"`

	VmallMemberRelationId uint64 `json:"vmall_member_relation_id" gorm:"not null;type:int(11);default:0;comment:VmallMemberRelation表id"`
	CustomerNameEnc       string `json:"customer_name" gorm:"not null;type:varchar(100);default:'';comment:客户名(加密)"`
	CustomerMobileEnc     string `json:"customer_mobile_enc" gorm:"not null;type:varchar(100);default:'';comment:客户手机号(加密)"`
	CustomerIdCodeEnc     string `json:"customer_id_code_enc" gorm:"not null;type:varchar(100);default:'';comment:客户身份证(加密)"`
	RelationType          int32  `json:"relation_type" gorm:"not null;type:int(11);default:0;comment:关系类型 1:亲戚 2:朋友"`
	StaffId               string `json:"staff_id" gorm:"not null;type:varchar(100);default:'';comment:员工id"`

	CreatedAt time.Time `json:"createdAt" gorm:"not null"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"not null"`
	DeletedAt time.Time `json:"deleted_at" gorm:"not null"`
}

type VmallMemberRelationExtandLog struct {
	Id uint64 `json:"id" structs:"id" gorm:"primary_key;autoIncrement:false"`

	VmallMemberRelationExtandId uint64 `json:"vmall_member_relation_id" gorm:"not null;type:int(11);default:0;comment:VmallMemberRelationExtand表id"`
	StaffId                     string `json:"staff_id" gorm:"not null;type:varchar(100);default:'';comment:员工id"`

	CreatedAt time.Time `json:"createdAt" gorm:"not null"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"not null"`
}

type EquityGoodsExtand struct {
	Id          uint64    `json:"id" gorm:"primary_key;autoIncrement:false"`
	GoodsId     uint      `json:"goods_id" gorm:"default:0;not null;comment:关联商品表id"`
	IsSyncVmall bool      `json:"is_sync_vmall" gorm:"default:0;not null;comment:是否同步到商城"`
	AdminId     uint64    `json:"admin_id" gorm:"default:0;not null;comment:操作人id"`
	CreatedAt   time.Time `json:"created_at"   gorm:"not null"`
	UpdatedAt   time.Time `json:"updated_at"   gorm:"not null"`
}

type EquitySupplier struct {
	gorm.Model
	Name                string            `gorm:"size:100; comment:供应商名称"`
	Image               string            `gorm:"size:100; comment:logo"`
	Introduction        string            `gorm:"type:text; comment:简介"`
	ServiceTel          string            `gorm:"size:100; comment:客服电话"`
	ManagerName         string            `gorm:"size:100; comment:负责人名字"`
	ManagerMobile       string            `gorm:"size:100; comment:负责人手机"`
	Attach              string            `gorm:"size:100; comment:内容用于将该供应商商品订单传到微信支付attach字段"`
	Remark              string            `gorm:"size:100; comment:备注"`
	ChannelID           string            `json:"channel_id" gorm:"not null;size:100;index"`
	RuiyinxinMerchantID string            `json:"ruiyinxin_merchant_id" gorm:"size:100; comment:第三方商户号（瑞银信）"`
	RuiyinxinChannel    string            `json:"ruiyinxin_channel" gorm:"size:100; comment:第三方通道（瑞银信）"`
	WyfMerchantID       string            `json:"wyf_merchant_id" gorm:"not null;type:varchar(100);comment:微邮付商户号"`
	NoWriteLimit        bool              `json:"no_write_limit" gorm:"not null;comment:放开限制(单商户核销每日限额)"`
	OrgID               string            `json:"org_id" gorm:"size:100; comment:组织架构ID"`
	CreateAdminID       string            `json:"create_admin_id" gorm:"not null;size:30; comment:创建的员工ID"`
	AttachmentImg       string            `json:"attachment_img" gorm:"type:varchar(100); not null; comment:附件图片"`
	ReviewStatus        uint32            `json:"review_status" gorm:"type:tinyint(4); not null; comment:审核状态：1-待确认 2-待审核 3-审核通过 4-审核不通过"`
	AssistEvidenceImg   stringArrayStruct `json:"assist_evidence_img" gorm:"type:text; comments:附件证明图片"`
	SignImg             string            `json:"sign_img" gorm:"type:varchar(200); not null; default:''; comments:签名图片"`
	ShelvesStatus       uint32            `json:"shelves_status" gorm:"type:tinyint(4); not null; default:1; comment:上下架状态：1-上架 2-下架"`
	BusinessType        int32             `json:"business_type" gorm:"type:tinyint(4); not null; default:0; comments:商户类型，1：企业；2：个体户；3：小微；4：其他；5：政府"`
}

func GetDb() *gorm.DB {
	return db
}

type EquityStockAmountLog struct {
	ID uint64 `json:"id" gorm:"primary_key;autoIncrement:false"`

	OrgID    uint64 `json:"org_id" gorm:"index:idx_org_goods;comment:网点ID"`
	GoodsID  uint64 `json:"goods_id" gorm:"index:idx_org_goods;comment:商品ID"`
	Amount   int64  `json:"amount" gorm:"comment:库存数量"`
	Method   string `json:"method" gorm:"type:varchar(100);comment:操作类型"`
	Source   string `json:"source" gorm:"type:varchar(100);comment:来源"`
	OrderNum string `json:"order_num" gorm:"type:varchar(50);default:'';comment:订单号"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

type ActivityPushMessage struct {
	Id         uint64    `json:"id" gorm:"primary_key;autoIncrement:false"`
	ReceiveId  uint64    `json:"receive_id" gorm:"index;not null;comment:领取记录Id"`
	Aid        uint64    `json:"aid" gorm:"index;not null;default:0;comment:活动id"`
	UserId     string    `json:"user_id" gorm:"index;not null;size:50;comment:用户id"`
	EventID    string    `json:"event_id" gorm:"not null;size:20;comment:事件Id"`
	EventType  string    `json:"eventType" gorm:"not null;size:20;comment:事件Id"`
	ActionId   string    `json:"actionid" gorm:"not null;size:100;comment:推送的活动参与流水号"`
	MobileNo   string    `json:"mobile_no" gorm:"index;not null;size:50;comment:手机号"`
	CstmName   string    `json:"cstm_name" gorm:"index;not null;size:50;comment:姓名"`
	IdCard     string    `json:"id_card" gorm:"index;not null;size:50;comment:身份证"`
	AwardSeq   string    `json:"award_seq" gorm:"not null;size:100;comment:推送的奖品流水号"`
	AwardName  string    `json:"award_name" gorm:"not null;size:50;comment:奖品名称"`
	AwardType  string    `json:"award_type" gorm:"not null;size:50;comment:奖品类型"`
	DoneDateAt time.Time `json:"done_date_at" gorm:"not null;comment:领取时间"`
	OrderBrch  string    `json:"order_brch" gorm:"not null;size:50;comment:机构号"`
	Status     int8      `json:"status" gorm:"not null;comment:推送状态：1、成功，0、待推送，-1、失败"`
	ErrMsg     string    `json:"err_msg" gorm:"not null;size:250;comment:推送失败原因"`
	ReqData    string    `json:"req_data" gorm:"not null;type:text;comment:推送请求"`
	ResData    string    `json:"res_data" gorm:"not null;type:text;comment:推送返回"`
	CreatedAt  time.Time `json:"createdAt" gorm:"not null"`
	UpdatedAt  time.Time `json:"updatedAt" gorm:"not null"`
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

type StringArrayStruct []string

func (array StringArrayStruct) Value() (driver.Value, error) {
	str := JsonMarshal(array)
	if str == "" {
		str = "[]"
	}
	return []byte(str), nil
}
func (array *StringArrayStruct) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	return json.Unmarshal(b, array)
}
