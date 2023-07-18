package testmodel

import "fmt"

func Exec() {
	goods := Goods{}
	id := 147128
	err := db.Where("id = ?", id).Take(&goods).Error
	//err := db.Debug().Where("name = ?", "月月邮好礼 鸡蛋").Take(&goods).Error

	fmt.Println(err)
	fmt.Println(goods.Id)

	//s := SysConfigInfo{}
	//db.Debug().Where("value = ?", "OFF").Take(&s)
	//fmt.Println(s)

}
