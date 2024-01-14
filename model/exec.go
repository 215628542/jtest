package model

import (
	"fmt"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"time"
)

func CreateTime(db *gorm.DB) *gorm.DB {
	return db.Where("created_at > ?", time.Time{})
}

func Exec() {

	orderIds := make([]int64, 0)
	db.Model(Tt{}).Where("id=1").Where("id=22").Pluck("id", &orderIds)
	//s.Weight = 99998877.12123
	//// 123123000
	//db.Model(&s).Where("id=486954396925038592").Save(&s)
	fmt.Println(len(orderIds))
	return
}

func Sum() {

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
