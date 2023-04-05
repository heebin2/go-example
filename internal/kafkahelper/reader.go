package kafkahelper

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
)

type Handler func(kafka.Message) error

type KafkaReader struct {
	Bootstrap string
	GroupID   string
	Topic     string
	handler   Handler
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
}

func NewKafkaReader(Bootstrap, GroupID, Topic string, handler Handler) (*KafkaReader, error) {
	if len(Bootstrap) <= 0 {
		return nil, fmt.Errorf("invalid Bootstrap (%s)", Bootstrap)
	}

	if len(GroupID) <= 0 {
		return nil, fmt.Errorf("invalid GroupID (%s)", GroupID)
	}

	if len(Topic) <= 0 {
		return nil, fmt.Errorf("invalid Topic (%s)", Topic)
	}

	if handler == nil {
		return nil, fmt.Errorf("invalid handler")
	}

	ctx, cancel := context.WithCancel(context.TODO())
	return &KafkaReader{
		Bootstrap: Bootstrap,
		GroupID:   GroupID,
		Topic:     Topic,
		handler:   handler,
		ctx:       ctx,
		cancel:    cancel,
	}, nil
}

func (kr *KafkaReader) Open() error {
	kr.wg.Add(1)
	go kr.run()

	return nil
}

func (kr *KafkaReader) Close() {
	kr.cancel()
	kr.wg.Wait()
}

func (kr *KafkaReader) newReader() *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:               strings.Split(kr.Bootstrap, ","),
		GroupID:               kr.GroupID,
		Topic:                 kr.Topic,
		MinBytes:              10e3,
		MaxBytes:              10e6,
		MaxWait:               1 * time.Second,
		ReadLagInterval:       -1,
		RebalanceTimeout:      1 * time.Second,
		ErrorLogger:           kafka.LoggerFunc(func(msg string, a ...any) { fmt.Printf("[KAFKA ERRO] "+msg, a...); fmt.Println() }),
		OffsetOutOfRangeError: true,
	})
}

func (kr *KafkaReader) run() {
	reader := kr.newReader()
	defer kr.wg.Done()
	defer reader.Close()

	for {
		msg, err := reader.ReadMessage(kr.ctx)
		if errors.Is(err, context.Canceled) {
			fmt.Println("ReadMessage context Canceled.")
			return
		}

		if err != nil {
			fmt.Println("ReadMessage fail : ", err)
			reader.Close()
			reader = kr.newReader()
			continue
		}

		if err := kr.handler(msg); err != nil {
			fmt.Println("handler fail : ", err)
		}
	}
}
