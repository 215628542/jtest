package main

import (
	"github.com/shopspring/decimal"
	"test/tool/redis"
)

func SwitchPrice(price int64) float64 {
	priceDec := decimal.NewFromInt(price)
	num := decimal.NewFromInt(1)
	dec := priceDec.Mul(num)
	f, _ := dec.Float64()
	return f
}

func main() {

	//antsTool.Test()
	//asynqTool.Run()
	//ants.Test()

	//testmodel.Exec()
	//fmt.Println("=======")
	//model.Exec()
	redis.Exec()
	//common.Differ()

	// 瑞银信签名
	//common.GenerateSignForData()

	// goodsInfo.GetData().GetSupplierDetail().GetRuiyinxinChannel()

}
