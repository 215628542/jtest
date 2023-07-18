package core

import (
	"encoding/json"
	"fmt"
	"runtime/debug"
	"time"
)

type GoWorker struct {
	// 协程池
	pool *AsyncQueue

	// 任务id
	TaskId int32

	// 协程空闲时间
	FreeTime int32
}

func (w *GoWorker) Run() {
	w.pool.addRunning(1)
	go func() {
		reqData := ""
		defer func() {
			w.pool.addRunning(-1)
			if perr := recover(); perr != nil {
				//fmt.Println("goworker-panic-reqData:" + reqData)
				if panicHandler := w.pool.options.PanicHandler; panicHandler != nil {
					// todo 将panic数据记录到mysql
					panicHandler(reqData, debug.Stack())

				} else {
					fmt.Printf("goworker exits from panic: %v\n%s\n", perr, debug.Stack())
				}
			}
			// 重新开启一个新的工作协程
			oldGoWorker := w.pool.workerArray[w.TaskId]
			goWorker := &GoWorker{
				pool:     w.pool,
				TaskId:   w.TaskId,
				FreeTime: 0,
			}
			w.pool.workerArray[w.TaskId] = goWorker
			goWorker.Run()
			// todo 如何正确销毁协程 销毁异常goWorker
			oldGoWorker = nil
			goWorker.pool.addRunning(1)
			fmt.Println("oldGoWorker====")
			fmt.Println(oldGoWorker)
		}()

		for {
			// 读取redis数据
			vals, err := w.pool.queue.get(w.TaskId)
			// 获取数据异常
			if err != nil {
				time.Sleep(1 * time.Second)
				continue
			}
			// 获取数据为空
			if len(vals) < 1 {
				// 获取空数据，说明没异步任务需要执行，空闲次数加1
				w.FreeTime++
				time.Sleep(1 * time.Second)
				continue
			}

			for _, val := range vals {
				fmt.Printf("获取异步消息%#v\n", val)
				reqData = val
				// 重置空闲次数
				w.FreeTime = 0
				funcErr := w.pool.options.TaskHandler(val)
				if funcErr != nil {
					fmt.Println("funcErr---" + funcErr.Error())
					// 如果任务返回err，则视为异常并进行重试
					asyncReqData := AsyncReqData{}
					jerr := json.Unmarshal([]byte(val), &asyncReqData)
					if jerr != nil {
						panic("解析请求数据异常:" + jerr.Error())
					}
					if asyncReqData.RetryNum > 0 {
						asyncReqData.RetryNum--
						w.pool.queue.Send(asyncReqData)
					}
				}
				w.pool.queue.DelTaskQueue(w.TaskId)
			}
		}
	}()

}
