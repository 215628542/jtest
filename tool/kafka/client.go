package kafka

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"sync"
)

const (
	GroupId = "first-group"
)

var Topics = []string{"activity"}

type Client struct {
	client         sarama.Client
	producer       sarama.AsyncProducer
	consumer       sarama.ConsumerGroup
	config         *Config
	driver         string
	consumeHandler []*consumeHandler
}

type Config struct {
	Servers  []string       `json:"servers"`
	Producer ProducerConfig `json:"producer"`
	Consumer ConsumerConfig `json:"consumer"`
}

type ProducerConfig struct {
	Has     bool `json:"has"`
	Retries int  `json:"retries"`
	Acks    int  `json:"acks"`
	Timeout int  `json:"timeout"`
}

type ConsumerConfig struct {
	Has     bool   `json:"has"`
	GroupId string `json:"groupId"`
	Oldest  bool   `json:"oldest"`
}

type consumeHandler struct {
	wg               *sync.WaitGroup
	cancel           context.CancelFunc
	setupFunc        func(session sarama.ConsumerGroupSession) error
	cleanupFunc      func(session sarama.ConsumerGroupSession) error
	consumeClaimFunc func(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error
}

var kafkaClient *Client

func GetKafkaClient() *Client {
	err := InitKafka()
	if err != nil {
		panic(err)
	}
	return kafkaClient
}

func InitKafka() error {
	if kafkaClient != nil && len(kafkaClient.config.Servers) > 0 {
		return nil
	}
	rc := &Config{
		Servers: []string{"192.168.1.105:9192"},
		Producer: ProducerConfig{
			Has:     true,
			Retries: 0,
			Acks:    0,
			Timeout: 0,
		},
		Consumer: ConsumerConfig{
			Has:     true,
			GroupId: GroupId,
			Oldest:  false,
		},
	}

	client, err := Connect(rc)
	if err != nil {
		return err
	}

	c := &Client{client: client, config: rc}
	if rc.Producer.Has {
		producer, _err := sarama.NewAsyncProducerFromClient(client)
		if _err != nil {
			return _err
		}
		c.producer = producer
		//定时刷新数据
		go func(p sarama.AsyncProducer) {
			errors := p.Errors()
			success := p.Successes()
			for {
				select {
				case err := <-errors:
					fmt.Println("kafka生产者异常-" + err.Error())
					break
				case msg := <-success:
					fmt.Println("kafka生产者成功")
					fmt.Println(msg)
					break
				}
			}
		}(c.producer)
	}
	if rc.Consumer.Has {
		consumer, _err := sarama.NewConsumerGroupFromClient(rc.Consumer.GroupId, client)
		if _err != nil {
			return _err
		}
		c.consumer = consumer
	}
	kafkaClient = c

	return nil
}

func Connect(c *Config) (sarama.Client, error) {
	sc := sarama.NewConfig()
	sc.Version = sarama.V1_1_0_0
	if c.Consumer.Oldest {
		fmt.Println("set oldest")
		sc.Consumer.Offsets.Initial = sarama.OffsetOldest
	} else {
		fmt.Println("not set oldest")
		sc.Consumer.Offsets.Initial = sarama.OffsetNewest
	}
	sc.Producer.Return.Successes = true
	sc.Producer.Return.Errors = true
	//sc.Net.MaxOpenRequests = c.MaxOpenRequest
	//sc.Net.DialTimeout = time.Duration(c.DialTimeout) * time.Second
	//sc.Net.ReadTimeout = time.Duration(c.ReadTimeout) * time.Second
	//sc.Net.WriteTimeout = time.Duration(c.WriteTimeout) * time.Second
	//sc.Admin.Timeout = time.Duration(c.Timeout) * time.Second
	//sc.Producer.MaxMessageBytes = c.MaxMessageBytes * 1024 * 1024
	//sc.Producer.Timeout = time.Duration(c.Timeout) * time.Second
	sc.Producer.Partitioner = sarama.NewRandomPartitioner
	return sarama.NewClient(c.Servers, sc)
}
