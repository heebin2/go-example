package main

import (
	"fmt"
	"go-helper/internal/kafkahelper"

	"github.com/segmentio/kafka-go"
)

const (
	Bootstrap = "localhost:9092,localhost:9093,localhost:9094"
)

func ExistTopic(bootstrap, topic string) bool {
	list, err := kafkahelper.TopicList(Bootstrap)
	if err != nil {
		return false
	}

	for i := range list {
		if list[i] == topic {
			return true
		}
	}

	return false
}

func main() {
	topicName := "temp"

	list, err := kafkahelper.TopicList(Bootstrap)
	if err != nil {
		fmt.Println("Topic list error : ", err)
		return
	}
	fmt.Println("List    : ", list)

	topicConfig := kafka.TopicConfig{
		Topic:             topicName,
		NumPartitions:     3,
		ReplicationFactor: 3,
	}
	if err := kafkahelper.TopicCreate(Bootstrap, topicConfig); err != nil {
		fmt.Println("create fail : ", err)
		return
	}

	if !ExistTopic(Bootstrap, topicName) {
		fmt.Println("create fail : not exist")
		return
	}
	fmt.Println("created :", topicName)

	if err := kafkahelper.TopicDelete(Bootstrap, topicName); err != nil {
		fmt.Println("delete fail : ", err)
		return
	}

	if ExistTopic(Bootstrap, topicName) {
		fmt.Println("delete fail : exist")
		return
	}
	fmt.Println("deleted :", topicName)
}
