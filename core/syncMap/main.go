package main

import (
	"fmt"
	"sync"
)

func main() {
	var syncMap sync.Map
	// 保存数据
	syncMap.Store("key1", "val1")
	syncMap.Store("key2", "val2")
	syncMap.Store("key3", "val3")

	// 删除数据
	//syncMap.Delete("key1")

	// 获取数据
	v, ok := syncMap.Load("key1")
	if !ok {
		panic("key1不存在")
	}

	v1, isString := v.(string)
	if !isString {
		panic("不是string类型")
	}
	fmt.Println(v1)

	syncMap.Range(func(key, value interface{}) bool {
		fmt.Printf("key:%s  value:%s \n", key, value)
		return key != "key2"
		//return true
	})

}
