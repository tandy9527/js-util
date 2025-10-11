package kafka

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	writer *kafka.Writer
}

var (
	producerOnce sync.Once
	producerInst *KafkaProducer
)

func InitProducer(brokers []string, topic string) {
	producerOnce.Do(func() {
		writer := &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Topic:        topic,
			Balancer:     &kafka.LeastBytes{},
			RequiredAcks: kafka.RequireAll,
			Async:        false,
		}
		producerInst = &KafkaProducer{writer: writer}
		log.Println("[Kafka] Producer initialized")
	})
}

func GetProducer() *KafkaProducer {
	if producerInst == nil {
		log.Fatal("[Kafka] Producer not initialized. Call InitProducer() first.")
	}
	return producerInst
}

func (p *KafkaProducer) SendMessage(key, value []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	msg := kafka.Message{
		Key:   key,
		Value: value,
		Time:  time.Now(),
	}
	return p.writer.WriteMessages(ctx, msg)
}

func (p *KafkaProducer) Close() {
	if p.writer != nil {
		_ = p.writer.Close()
		log.Println("[Kafka] Producer closed")
	}
}
