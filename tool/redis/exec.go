package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/cast"
	"time"
)

func Exec() {

	r, e := rdb.Get(context.Background(), "aaa").Result()
	fmt.Println(r)
	fmt.Println("=====")
	fmt.Println(e)
	if e != nil {
		fmt.Println("===aabbcc==")
	}
	return

	TestOrderTaskPool()
	return

	key := "test_list"
	ctx := context.Background()
	d, err := rdb.LPop(ctx, key).Result()

	if err != nil {
		fmt.Println("123")
	}

	fmt.Println(err)
	fmt.Println(d)

	//ZRangeByScore()
	//key := "t1"
	//ctx := context.Background()
	//rdb.ZAdd()

	//Ping()
	//appendSlice()
}

func TestOrderTaskPool() {

	type RequestBody struct {
		Unique string `json:"unique"`
		Data   string `json:"data"`
	}
	type RequestData struct {
		Key             string
		RequestId       int64
		HandlerFuncType string
		Data            []byte
	}

	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
	for i := 1; i <= 1; i++ {

		requestBody := RequestBody{
			Unique: cast.ToString(i),
			Data:   "",
		}
		requestBodyByte, _ := json.Marshal(requestBody)
		req := &RequestData{
			Data:            requestBodyByte,
			HandlerFuncType: "wx_pay_notify",
		}
		reqTmp, _ := json.Marshal(req)

		ctx := context.Background()
		cacheKey := "vmallorder:taskPool:queue"
		rdb.RPush(ctx, cacheKey, string(reqTmp))
	}
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
}

// 有序集合写入
func ZAdd() {
	members := &redis.Z{
		Score:  cast.ToFloat64(time.Now().Unix()),
		Member: "bb",
	}
	id, err := rdb.ZAdd(context.Background(), "jh", members).Result()
	fmt.Println(id) // 第一次写入返回1，同一个member再次写入返回0
	fmt.Println(err)
}

// 读取有序集合
func ZRangeByScore() {
	val, err := rdb.ZRangeByScore(context.Background(), "jh", &redis.ZRangeBy{
		Min:    "-inf",
		Max:    "1689496063",
		Offset: 0,
		Count:  1,
	}).Result()
	fmt.Println(err)

	// 删除集合元素
	id, err := rdb.ZRem(context.Background(), "jh", val).Result()

	fmt.Println(val)
	fmt.Println(id)
	fmt.Println(err)
}

func GetSet() {
	key := "append"
	ctx := context.Background()

	memberIds := []string{"3", "4"}
	cacheData, err := rdb.Get(ctx, key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		panic(err)
	}
	mids := make([]string, 0)
	json.Unmarshal([]byte(cacheData), &mids)

	mids = append(mids, memberIds...)
	mids = removeRepByMap(mids)
	fmt.Println(mids)

	memberIdStr, err := json.Marshal(mids)
	if err != nil {
		panic(err)
	}
	rdb.Set(ctx, key, string(memberIdStr), time.Hour*24*30)

}

// slice去重
func removeRepByMap(slc []string) []string {
	result := []string{}
	tempMap := map[string]byte{}
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e] = 0

		if len(tempMap) != l {
			result = append(result, e)
		}
	}
	return result
}

func lock() {

	ctx := context.Background()
	locker := NewRedisLocker(ctx, rdb, "equitygoods-crontab-writeofffailtry1", 1*time.Hour)

	b := locker.Lock()
	fmt.Println(b)

	if b {
		// 解🔐
		defer locker.Unlock()

		fmt.Println("111")
	}
	sErr := rdb.Set(ctx, "aa", 22, 1*time.Minute).Err()
	fmt.Println(sErr)

}

func Ping() {
	pong, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println("reis 连接失败：", pong, err)
		panic("reis 连接失败")
	}
	fmt.Println("连接成功")
}
