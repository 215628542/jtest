package kafka

import (
	"encoding/json"
	"github.com/IBM/sarama"
)

func GetProducer() sarama.AsyncProducer {
	kc := GetKafkaClient()
	return kc.producer
}

func SendMsg(key, value string) {
	msgBts, _ := json.Marshal(value)
	msg := &sarama.ProducerMessage{
		Topic: Topics[0],
		Key:   sarama.ByteEncoder(key),
		Value: sarama.ByteEncoder(msgBts),
	}
	GetProducer().Input() <- msg
}
