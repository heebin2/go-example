package main

import (
	"fmt"
	"go-helper/internal/kafkahelper"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	Bootstrap = "localhost:9092,localhost:9093,localhost:9094"
	Topic     = "test"
)

func main() {
	writer, err := kafkahelper.NewKafkaWriter(Bootstrap, Topic)
	if err != nil {
		fmt.Println("NewKafkaWriter fail : ", err)
		return
	}

	if err := writer.Open(); err != nil {
		fmt.Println("kafak writer open fail : ", err)
		return
	}
	defer writer.Close()

	ticker := time.NewTicker(100 * time.Millisecond)
	ticker2 := time.NewTicker(10 * time.Second)
	for {
		select {
		case now := <-ticker.C:
			writer.In <- kafka.Message{Value: []byte(now.String())}
		case <-ticker2.C:
			fmt.Println("Exit try main")
			return
		}
	}
}
