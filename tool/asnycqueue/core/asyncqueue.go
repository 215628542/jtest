package core

import (
	"errors"
	"fmt"
	"github.com/gofrs/uuid"
	"sync/atomic"
	"test/tool/redis"
)

var AsyncQueues map[string]*AsyncQueue

type AsyncQueue struct {

	// 配置信息
	options *Options

	// 正在运行协程数
	running int32

	// 消息队列
	queue Queue

	// 协程组
	workerArray []*GoWorker
}

func NewAsyncQueue(options ...Option) (*AsyncQueue, error) {
	opts := loadOptions(options...)

	// 检查配置信息
	if opts.AsyncQueueName == "" {
		return nil, errors.New("请配置异步消息任务名")
	}
	if opts.PoolSize < 1 {
		return nil, errors.New("请配置协程池大小")
	}
	if opts.TaskHandler == nil {
		return nil, errors.New("请配置消费任务")
	}

	pool, ok := AsyncQueues[opts.AsyncQueueName]
	if ok {
		return pool, nil
	}

	p := &AsyncQueue{
		options: opts,
		queue: Queue{
			Name: opts.AsyncQueueName,
			Rdb:  redis.GetRdb(),
		},
	}

	// 启动worker协程
	p.workerArray = make([]*GoWorker, opts.PoolSize)
	var taskId int32 = 0
	for ; taskId < opts.PoolSize; taskId++ {
		// 启动worker协程
		goWorker := &GoWorker{
			pool:     p,
			TaskId:   taskId,
			FreeTime: 0,
		}
		p.workerArray[taskId] = goWorker
		goWorker.Run()
	}
	fmt.Printf("初始化启动的协程数%d\n", len(p.workerArray))

	return p, nil
}

// 发送异步任务消息
func (a *AsyncQueue) EnQueue(reqData AsyncReqData) error {
	reqData.id = uuid.Must(uuid.NewV4()).String()
	return a.queue.Send(reqData)
}

func (a *AsyncQueue) addRunning(delta int) {
	atomic.AddInt32(&a.running, int32(delta))
}

func (a *AsyncQueue) FmtRunning() {
	fmt.Printf("当前running数值--%d\n", a.running)
}

func (a *AsyncQueue) FmtWorkArray() {
	fmt.Println(a.workerArray)
	for k, v := range a.workerArray {
		if v != nil {
			fmt.Printf("当前workerArray数值--ptaskId:%d-taskId:%d\n", k, v.TaskId)
		} else {
			fmt.Println("协程池中的协程有异常")
		}
	}

}
