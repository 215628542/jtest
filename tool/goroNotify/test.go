package main

import (
	"context"
	"fmt"
	"test/model"
	"time"
)

//var db = model.GetDb()
//var rdb = redis.GetRdb()

func main1() {

	// 1、进来加锁，忽略掉 NewMateriallibraryClient
	// 2、根据id清空数据，清空redis数据
	// 3、生成taskId，并且生成redis数据
	// 4、两个协程之间进行通知

	ctx := context.Background()
	id := 1
	var taskId int64 = 1
	key := "task_id_%d"
	taskKey := fmt.Sprintf(key, taskId)

	// 获取需要删除的taskId
	delList := []model.Tt{}
	db.Where("id=?", id).Find(&delList)

	// 2、根据id清空数据
	record := model.Tt{}
	db.Where("id=?", id).Delete(&record)

	for _, t := range delList {
		// 2、清空redis数据
		delkKey := fmt.Sprintf(key, t.TaskId)
		rdb.Del(ctx, delkKey)
	}
	// 3、生成taskId，并且生成redis数据
	rdb.Set(ctx, taskKey, 1, 3600*time.Second)

	go do(taskId, taskKey)

	time.Sleep(2000 * time.Second)
}

func do(taskId int64, taskKey string) {
	stopRecordChan := make(chan struct{}, 1)
	stopListenChan := make(chan struct{}, 1)

	// 使用ctx进行通知任务过期，停止任务继续执行
	recordCtx := context.Background()
	// 生成数据
	go makeRecord(recordCtx, taskId, stopRecordChan, stopListenChan)
	// 监听任务是否被取消
	go listen(recordCtx, taskKey, stopRecordChan, stopListenChan)

	//select {
	//case <-stopRecordChan:
	//	// 删除任务生成的数据
	//case <-stopListenChan:
	//}

}

func makeRecord(ctx context.Context, taskId int64, stopRecordChan chan struct{}, stopListenChan chan struct{}) {

	fmt.Println("en makeRecord")
	for i := 0; i < 300; i++ {

		select {
		case <-stopRecordChan:
			fmt.Println("channel stop")
			db.Where("task_id=?", taskId).Delete(model.Tt{})
			return
		default:
			fmt.Println("make record")
			time.Sleep(3 * time.Second)
			t := model.Tt{Id: 1, TaskId: taskId, CreateTime: time.Now()}
			db.Create(&t)
		}
	}

	// 数据生成完成，则通知监控协程退出
	stopListenChan <- struct{}{}
	fmt.Println("生成数据完成")
}

func listen(ctx context.Context, taskKey string, stopRecordChan chan struct{}, stopListenChan chan struct{}) {

	// 每秒读取一次
	for {
		select {
		case <-time.Tick(5 * time.Second):
			fmt.Println("监控任务是否被取消")
			isExist := rdb.Exists(context.Background(), taskKey).Val()
			fmt.Println(taskKey)
			fmt.Println(isExist)
			if isExist < 1 {
				stopRecordChan <- struct{}{}
				return
			}
		case <-stopListenChan:
			fmt.Println("已停止监控协程")
			return
		}
	}
}
