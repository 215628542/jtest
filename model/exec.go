package model

import (
	"context"
	"fmt"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

func CreateTime(db *gorm.DB) *gorm.DB {
	return db.Where("created_at > ?", time.Time{})
}

func Exec() {

	//t := Tt{}
	var cnt1 int64
	ctx := context.Background()
	db.WithContext(ctx).Model(Tt{}).Select("count(distinct task_id) as cnt ").
		Where(" id > ? ", 0).
		//Group("task_id").
		Order("id asc").
		Count(&cnt1)
	fmt.Println(cnt1)

	return

	tx := db.WithContext(ctx).Model(Activity{}).
		Where("deleted_at=? and channel_id=? and status = ? ", time.Time{}, "123", 1).
		Limit(5).
		Order("sort desc")

	// 匹配活动有效范围内
	tx = tx.Where("start_date < ? and end_date > ?", time.Now().Format("2006-01-02 15:04:05"), time.Now().Format("2006-01-02 15:04:05"))
	// 过滤可推广人群
	tx = tx.Where(" find_in_set(?, promo_crowds) ", "merchant")

	fmt.Println("PromoActivityPage-test4")
	list := []ActivityBonus{}
	tx.Find(&list)

	var total int64
	err := db.WithContext(ctx).Model(&ActivityBonus{}).
		Select("COALESCE(SUM(bonus),0)").
		//Where(" promo_member_id = ? and deleted_at=? and bonus>0 and channel_id=? and status=?", "1", time.Time{}, "1", "1").
		Scan(&total).Error

	fmt.Println(err)
	fmt.Println(total)

	//var totalAmount float64
	//db.Model(&Order{}).Where("status = ?", "completed").Select("SUM(amount)").Scan(&totalAmount)

	return

	//tt := Tt{}
	//tx := db.Begin()
	//tx.Set("gorm:query_option", "FOR UPDATE").First(&tt, "id")
	//fmt.Println(tt)
	//return

	type Gg struct {
		Num int
	}

	g := []Gg{}
	db.Table("tts").Model(Tt{}).Select("count(id) as num").Where("id>?", 0).
		Group("remark").Order("num desc").Scan(&g)
	fmt.Println(g)

	return

	for {
		time.Sleep(time.Second)

		rand.Seed(time.Now().UnixNano())
		id := rand.Intn(4)

		if id < 1 {
			continue
		}

		t2 := Tt{Id: cast.ToUint32(id)}

		gid, err := t2.Get(ctx, db, id)
		if err != nil {
			fmt.Println(cast.ToUint32(id))
			panic(err)
		}
		fmt.Println(gid)
		//fmt.Println(t2.Remark)
	}

	orderIds := make([]int64, 0)
	w := "12"
	db.Model(Tt{}).Where("remark = ?", w).Pluck("id", &orderIds)
	//s.Weight = 99998877.12123
	//// 123123000
	//db.Model(&s).Where("id=486954396925038592").Save(&s)
	fmt.Println(orderIds)

	return
}

func Sum2() {

	rec := &struct {
		Snum int
	}{}

	db.Model(&WriteOffFail{}).Unscoped().Select("SUM(try_num) as snum").Where("id>1").Take(rec)
	fmt.Println(rec)
}

// 获取单列多行
func GetColms() {
	ids := make([]string, 0)
	db.Model(&WriteOffFail{}).Unscoped().Where("id=1").Pluck("id", &ids)
	fmt.Println(ids)
}

func CreateData() {
	createTime, _ := time.ParseInLocation("2006-01-02 15:04:05", "1906-06-06 06:06:06", time.Local)
	model := VmallOrderPay{
		Id:          cast.ToUint64(time.Now().Unix()),
		SuccessedAt: time.Now(),
		CreatedAt:   createTime,
		UpdatedAt:   time.Now(),
	}
	err := db.Create(&model).Error
	if err != nil {
		panic(err)
	}
}

func Sum() {
	var total int64
	db.WithContext(context.Background()).Model(&ActivityBonus{}).
		Select("COALESCE(SUM(bonus),0)").
		//Where(" promo_member_id = ? and deleted_at=? and bonus>0 and channel_id=? and status=?", "1", time.Time{}, "1", "1").
		Scan(&total)
}
