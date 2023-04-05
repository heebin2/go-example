package kafkahelper

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
)

func TopicList(bootstrap string) ([]string, error) {
	s := strings.Split(bootstrap, ",")
	if len(s) < 1 {
		return nil, fmt.Errorf("invalid bootstrap")
	}

	conn, err := kafka.Dial("tcp", s[0])
	if err != nil {
		return nil, errors.Wrap(err, "kafka.Dial fail")
	}
	defer conn.Close()

	partitions, err := conn.ReadPartitions()
	if err != nil {
		return nil, errors.Wrap(err, "ReadPartitions fail")
	}

	ret := []string{}
	retemp := make(map[string]any)
	for i := range partitions {
		if partitions[i].Topic != "__consumer_offsets" {
			retemp[partitions[i].Topic] = true
		}
	}

	for k := range retemp {
		ret = append(ret, k)
	}

	return ret, nil
}

func TopicCreate(bootstrap string, topic ...kafka.TopicConfig) error {
	s := strings.Split(bootstrap, ",")
	if len(s) < 1 {
		return fmt.Errorf("invalid bootstrap")
	}

	conn, err := kafka.Dial("tcp", s[0])
	if err != nil {
		return errors.Wrap(err, "kafka.Dial fail")
	}
	defer conn.Close()

	return conn.CreateTopics(topic...)
}

func TopicDelete(bootstrap string, topic ...string) error {
	s := strings.Split(bootstrap, ",")
	if len(s) < 1 {
		return fmt.Errorf("invalid bootstrap")
	}

	conn, err := kafka.Dial("tcp", s[0])
	if err != nil {
		return errors.Wrap(err, "kafka.Dial fail")
	}
	defer conn.Close()

	return conn.DeleteTopics(topic...)
}
