package main

import (
	"fmt"
	"go-helper/internal/kafkahelper"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	Bootstrap = "localhost:9092,localhost:9093,localhost:9094"
	GroupID   = "test"
	Topic     = "tms.result"
)

func MessageHandler(msg kafka.Message) error {
	if l := len(msg.Value); l >= 16 {
		fmt.Println("Topic=", msg.Topic, ", Message Length=", l, ", Message Key=", msg.Key, ", Message Header=", msg.Headers)
		return nil
	}

	fmt.Println("Topic=", msg.Topic, ",Message Length=", msg.Value, ", Message Key=", msg.Key, ", Message Header=", msg.Headers)
	return nil
}

func main() {
	reader, err := kafkahelper.NewKafkaReader(Bootstrap, GroupID, Topic, MessageHandler)
	if err != nil {
		fmt.Println("NewKafkaReader fail :", err)
		return
	}

	if err := reader.Open(); err != nil {
		fmt.Println("kafka reader open fail : ", err)
		return
	}
	defer reader.Close()

	time.Sleep(10 * time.Second)
	fmt.Println("Exit try main")
}
