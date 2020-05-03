package kafka

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Shopify/sarama"

	"github.com/alka/supermart/store/model"
)

const (
	kafkaConn = "localhost:9092"
	topic     = "supermart"
)

//NewProducer create new producer
func NewProducer() (sarama.SyncProducer, error) {
	// setup sarama log to stdout
	sarama.Logger = log.New(os.Stdout, "", log.Ltime)

	// producer config
	config := sarama.NewConfig()
	config.Producer.Retry.Max = 5
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true

	// sync producer
	prd, err := sarama.NewSyncProducer([]string{kafkaConn}, config)

	return prd, err
}

//PublishMsg publish message to kafka
func PublishMsg(items []*model.Item, producer sarama.SyncProducer) {
	var partition int32
	var offset int64

	for _, i := range items {
		message, err := json.Marshal(i)
		msg := &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.StringEncoder(string(message)),
		}
		partition, offset, err = producer.SendMessage(msg)
		if err != nil {
			fmt.Println("Error publish: ", err.Error())
		}
	}
	// publish sync

	fmt.Println("Partition: ", partition)
	fmt.Println("Offset: ", offset)
}
