package model

import (
	"context"
	"fmt"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"time"
)

func CreateTime(db *gorm.DB) *gorm.DB {
	return db.Where("created_at > ?", time.Time{})
}

func Exec() {

	list := make([]ActivityData, 0)
	err := db.WithContext(context.Background()).Table(ActivityData{}.TableName()).
		Where("cid = ? ", 432492514382188544).Find(&list).Error
	fmt.Println(err)
	fmt.Println(list)
	return

	//t := VmallMemberPoint{Mobile: "123", ChannelId: "456"}
	//err := db.Save(&t).Error
	//fmt.Println(err)
	//fmt.Println(t)
	//errors.Is(err, gorm.d)
	//
	//return

	tx := db.Model(VmallInterPointCheckIn{})

	data := VmallInterPointCheckIn{}
	err = tx.First(&data).Error

	// json_contains(limit_org_types,json_array("%s"))

	fmt.Println(err)
	//fmt.Println(data)
	fmt.Printf("%#v", data)
	//fmt.Printf("%#v", data.CheckInType[1])

	return
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
