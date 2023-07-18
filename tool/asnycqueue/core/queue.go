package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/shopspring/decimal"
	"strconv"
	"time"
)

const (
	QueueTest = "testQueue"
)

type Queue struct {

	// 消息队列名称
	Name string

	// 请求数据
	//ReqData AsyncReqData

	// redis缓存
	Rdb *redis.Client
}

type AsyncReqData struct {
	// 延迟执行时间
	DelayTime time.Duration

	// 重试次数
	RetryNum int32

	// 请求数据
	Value string

	// id
	id string
}

// 从队列获取消息消费
func (q *Queue) get(taskId int32) (vals []string, err error) {
	ctx := context.Background()
	// 获取当前时间戳内的消息消费
	vals, err = q.Rdb.ZRangeByScore(ctx, q.Name, &redis.ZRangeBy{
		Min:    "-inf",
		Max:    strconv.FormatInt(time.Now().Unix(), 10),
		Offset: 0,
		Count:  1,
	}).Result()
	if err != nil {
		return
	}

	id, err := q.Rdb.ZRem(ctx, q.Name, vals).Result()
	if err != nil {
		err = errors.New("移除消息队列数据失败")
		return
	}
	if id < 1 {
		err = errors.New("删除了空的集合数据，说明消息已经被其他协程抢占了")
		return
	}

	// 新增redis hash，reqData新增id，goworker消费完，再根据id删除hash数据
	for _, val := range vals {
		// 将消息放入redis hash，和worker对应的消息渠道中
		q.Rdb.HSet(ctx, q.QueueKey(), q.QueueTaskKey(taskId), val)
	}

	// todo 使用lua将上述redis封装成原子操作

	return
}

// 发送消息到消息队列
func (q *Queue) Send(reqData AsyncReqData) error {

	ctx := context.Background()
	msg, err := json.Marshal(reqData)
	if err != nil {
		return err
	}
	members := &redis.Z{
		Score:  q.Score(ctx, time.Now().Add(reqData.DelayTime).Unix()),
		Member: msg,
	}
	_, err = q.Rdb.ZAdd(ctx, q.Name, members).Result()
	return err
}

// goworker执行任务完成后，删除任务消息
func (q *Queue) DelTaskQueue(taskId int32) error {
	_, err := q.Rdb.HDel(context.Background(), q.QueueKey(), q.QueueTaskKey(taskId)).Result()
	return err
}

func (q *Queue) QueueKey() string {
	return fmt.Sprintf("asyncQueue-queueName-%s", q.Name)
}

func (q *Queue) QueueTaskKey(taskId int32) string {
	return fmt.Sprintf("asyncQueue-queueName-%s-taskId-%d", q.Name, taskId)
}

func (q *Queue) Score(ctx context.Context, t int64) float64 {
	priceDec := decimal.NewFromInt(t)
	num := decimal.NewFromInt(1)
	dec := priceDec.Mul(num)
	f, _ := dec.Float64()
	return f
}
