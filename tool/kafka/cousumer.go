package kafka

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"runtime/debug"
)

func InitCousumer() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("kafka消费者启动异常") // 打印错误记录
			fmt.Println(err)
			fmt.Println(string(debug.Stack())) // 记录堆栈
		}
	}()
	kc := GetKafkaClient()

	// 创建消费者组
	cg, err := sarama.NewConsumerGroupFromClient(GroupId, kc.client)
	if err != nil {
		panic("创建消费者失败：" + err.Error())
	}
	defer cg.Close()

	// 阻塞运行
	consume(&cg, Topics, GroupId)
}

func consume(group *sarama.ConsumerGroup, topics []string, groupId string) {
	fmt.Println("消费者组--" + groupId + "--开始监听")
	ctx := context.Background()
	for {
		handler := consumerGroupHandler{name: groupId}
		err := (*group).Consume(ctx, topics, handler)
		if err != nil {
			panic(err)
		}
	}
}

type consumerGroupHandler struct {
	name string
}

func (consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error { return nil }

func (consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {

		//ConsumeTest(sess, msg)
		fmt.Println(string(msg.Key))
		fmt.Println(string(msg.Value))
		fmt.Println("================")
		switch string(msg.Key) {
		case "testKey":
			fmt.Println(string(msg.Key))
			fmt.Println(string(msg.Value))
			// 确认消息
			sess.MarkMessage(msg, "")

		default:
			//sess.MarkMessage(msg, "")
			//fmt.Println("default makeMessage")
		}
	}
	return nil
}

func ConsumeTest(sess sarama.ConsumerGroupSession, msg *sarama.ConsumerMessage) {
	// 确认消息
	defer sess.MarkMessage(msg, "")

	fmt.Println("ConsumeTest-=================")
	fmt.Println(string(msg.Key))
	fmt.Println(string(msg.Value))
}
