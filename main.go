package main

import (
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
	"test/core/shangyi"
	"test/tool/kafka"
	"time"
)

func SwitchPrice(price string) (float64, error) {
	priceDec, err := decimal.NewFromString(price)
	if err != nil {
		return 0, err
	}
	num := decimal.NewFromInt(100)
	dec := priceDec.Mul(num)
	res, _ := dec.Float64()
	return res, err
}

func main() {
	kafka.RunKafka()

	//kafka.SendMsg("testKey", "testt")

	time.Sleep(time.Hour)
	return

	//model.Exec()
	return

	//decrypt()

	//model.Exec()
	//model.Test()
	//model.Exec()
	//model.Test()
	//model.Test()
	return

	//cstZone := time.FixedZone("CST", 8*3600) // 东八
	//return time.Unix(t, 0).In(time.FixedZone("CST", 8*3600))

	if time.Now().Before(cast.ToTime("2024-06-13 15:54:01")) {
		fmt.Println(123123)
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
		fmt.Println(cast.ToTime("2024-06-13 15:54:01").Format("2006-01-02 15:04:05"))
	}
	if time.Now().Format("2006-01-02 15:04:05") <= "2024-06-13 16:03:01" {
		fmt.Println("------")
	}
	return

	// 解密
	//decrypt()

	// 执行数据库
	//model.Exec()

	//elasticSearch.Exec()

	//antsTool.Test()
	//asynqTool.Run()
	//ants.Test()

	//fmt.Println("=======")
	//model.Exec()
	//redis.Exec()
	//common.Differ()

	// 瑞银信签名
	//common.GenerateSignForData()

	// goodsInfo.GetData().GetSupplierDetail().GetRuiyinxinChannel()

}

// 解密
func decrypt() {
	s :=

		"9Cg81Tt+EZo6MV72V3jDXnaaiKWS9R6+kyDBo4HiyuecvM2X2wo5K7oFGtMFusdWNLp0vnBFX3F+WWbOcRT3D7NkfLJpK9Z5p9hrninNe/fqg0lHOcw0EC9H0/OVB12U"
	//"VxTAOc+db/Mx+io/cH+cFTE65ejgQcx+eWO4IkgifLH/sUqCeP61kkpbFl6UzSe0yjuXFPKEuxdhAOoSz8knukFWn4WeOSR2MLOzwBLm8Yn0O9o4oq2bBm2O16nLl9hhamj2AZyYYw+i2ygB2e4ECLUHlVhk7C2j9NoSDNm9bM+6U6z2+aqY5sS7iSuro6aDm2zxeLwtj/h8KWybkHe9t1ftDcxNa2V5+8+PpkGXTQY="
	a, e := shangyi.Decrypt(s, "company")
	fmt.Println(e)
	fmt.Println(a)
}
