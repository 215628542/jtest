package main

import (
	"fmt"
	"github.com/thoas/go-funk"
	"time"
)

// 切片合并
var dealIds = func(combineIds *[]string, join *[]string) int {
	if len(*combineIds) == 0 {
		*combineIds = *join
	} else {
		*combineIds = funk.InnerJoinString(*combineIds, *join) // 取交集
	}
	return len(*combineIds)
}

func main() {

	f := time.Now().Format("2006-01-02")
	fmt.Println(f)

	return

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
