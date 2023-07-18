package main

import (
	"encoding/json"
	"fmt"
	"runtime"
	"test/tool/asnycqueue/core"
	"test/tool/asnycqueue/taskFunc"
	"time"
)

func main() {
	pool, err := core.NewAsyncQueue(core.WithAsyncQueueName("test"),
		core.WithPoolSize(3),
		core.WithTaskHandler(taskFunc.Test),
		core.WithPanicHandler(taskFunc.TestPanicFunc),
	)

	msg := taskFunc.TestQueueMsg{
		TestMsg: "测试输入数据---1111",
	}
	m, _ := json.Marshal(msg)
	pool.EnQueue(core.AsyncReqData{
		DelayTime: 2 * time.Second,
		RetryNum:  2,
		Value:     string(m),
	})

	msg2 := taskFunc.TestQueueMsg{
		TestMsg: "测试输入数据---2222",
	}
	m2, _ := json.Marshal(msg2)
	pool.EnQueue(core.AsyncReqData{
		DelayTime: 2 * time.Second,
		RetryNum:  2,
		Value:     string(m2),
	})

	if err != nil {
		panic(err)
	}

	for {
		select {
		case <-time.Tick(3 * time.Second):
			fmt.Printf("当前运行的协程数：%d\n", runtime.NumGoroutine())
			pool.FmtRunning()
			pool.FmtWorkArray()
		}
	}

	time.Sleep(10 * time.Minute)
}
