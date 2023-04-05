package kafkahelper

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/snappy"
)

type KafkaWriter struct {
	Bootstrap string
	Topic     string
	In        chan kafka.Message
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
}

func NewKafkaWriter(Bootstrap, Topic string) (*KafkaWriter, error) {
	if len(Bootstrap) <= 0 {
		return nil, fmt.Errorf("invalid Bootstrap (%s)", Bootstrap)
	}

	if len(Topic) <= 0 {
		return nil, fmt.Errorf("invalid Topic (%s)", Topic)
	}

	ctx, cancel := context.WithCancel(context.TODO())
	return &KafkaWriter{
		Bootstrap: Bootstrap,
		Topic:     Topic,
		In:        make(chan kafka.Message),
		ctx:       ctx,
		cancel:    cancel,
	}, nil
}

func (kw *KafkaWriter) Open() error {
	kw.wg.Add(1)
	go kw.run()

	return nil
}

func (kw *KafkaWriter) Close() {
	kw.cancel()
	kw.wg.Wait()
}

func (kw *KafkaWriter) newWriter() *kafka.Writer {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:          strings.Split(kw.Bootstrap, ","),
		Topic:            kw.Topic,
		BatchSize:        100,
		BatchBytes:       1000 * 1000 * 100,
		BatchTimeout:     1 * time.Second,
		CompressionCodec: snappy.NewCompressionCodec(),
		ErrorLogger:      kafka.LoggerFunc(func(msg string, a ...any) { fmt.Printf("[KAFKA] "+msg, a...); fmt.Println() }),
		Async:            true,
		ReadTimeout:      1 * time.Second,
	})
	writer.AllowAutoTopicCreation = true

	return writer
}

func (kw *KafkaWriter) run() {
	writer := kw.newWriter()
	defer close(kw.In)
	defer kw.wg.Done()
	defer writer.Close()

	for {
		select {
		case <-kw.ctx.Done():
			return
		case msg := <-kw.In:
			err := writer.WriteMessages(kw.ctx, msg)
			if errors.Is(err, context.Canceled) {
				return
			}

			if err != nil {
				fmt.Println("WriteMessage fail : ", err)
				writer.Close()
				writer = kw.newWriter()
				continue
			}

			fmt.Println("write success : ", string(msg.Value))
		}
	}
}
