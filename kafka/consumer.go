package kafka

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

type MessageHandler func(ctx context.Context, key, value []byte) error

type KafkaConsumer struct {
	readers     []*kafka.Reader
	handler     MessageHandler
	ctx         context.Context
	cancel      context.CancelFunc
	wg          sync.WaitGroup
	retryCount  int
	retryDelay  time.Duration
	dlqTopic    string
	brokers     []string
	dlqProducer *KafkaProducer
}

var (
	consumerOnce sync.Once
	consumerInst *KafkaConsumer
)

func InitConsumer(cfg Config, handler MessageHandler, retryCount int, retryDelay time.Duration, dlqTopic string) {
	consumerOnce.Do(func() {
		ctx, cancel := context.WithCancel(context.Background())
		c := &KafkaConsumer{
			handler:    handler,
			ctx:        ctx,
			cancel:     cancel,
			retryCount: retryCount,
			retryDelay: retryDelay,
			dlqTopic:   dlqTopic,
			brokers:    cfg.Brokers,
		}

		if dlqTopic != "" {
			c.dlqProducer = &KafkaProducer{writer: &kafka.Writer{
				Addr:     kafka.TCP(cfg.Brokers...),
				Topic:    dlqTopic,
				Balancer: &kafka.LeastBytes{},
			}}
		}

		for _, topic := range cfg.Topics {
			r := kafka.NewReader(kafka.ReaderConfig{
				Brokers:        cfg.Brokers,
				GroupID:        cfg.GroupID,
				Topic:          topic,
				MinBytes:       cfg.MinBytes,
				MaxBytes:       cfg.MaxBytes,
				CommitInterval: time.Duration(cfg.CommitInterval) * time.Millisecond,
			})
			c.readers = append(c.readers, r)
		}

		consumerInst = c
		go consumerInst.start()
		log.Println("[Kafka] Consumer initialized with retry and DLQ")
	})
}

func GetConsumer() *KafkaConsumer {
	if consumerInst == nil {
		log.Fatal("[Kafka] Consumer not initialized. Call InitConsumer() first.")
	}
	return consumerInst
}

func (c *KafkaConsumer) start() {
	log.Println("[Kafka] Consumer started")
	c.wg.Add(1)
	defer c.wg.Done()

	for _, reader := range c.readers {
		go func(r *kafka.Reader) {
			for {
				m, err := r.FetchMessage(c.ctx)
				if err != nil {
					if err == context.Canceled {
						return
					}
					log.Printf("[Kafka] fetch error: %v", err)
					time.Sleep(time.Second)
					continue
				}

				success := false
				for attempt := 0; attempt <= c.retryCount; attempt++ {
					if err := c.handler(c.ctx, m.Key, m.Value); err != nil {
						log.Printf("[Kafka] handler error, attempt %d: %v", attempt+1, err)
						time.Sleep(c.retryDelay)
					} else {
						success = true
						break
					}
				}

				if !success && c.dlqProducer != nil {
					_ = c.dlqProducer.SendMessage(m.Key, m.Value)
					log.Printf("[Kafka] message sent to DLQ: key=%s", string(m.Key))
				}

				if err := r.CommitMessages(c.ctx, m); err != nil {
					log.Printf("[Kafka] commit error: %v", err)
				}
			}
		}(reader)
	}
}

func (c *KafkaConsumer) Stop() {
	c.cancel()
	c.wg.Wait()
	for _, r := range c.readers {
		_ = r.Close()
	}
	if c.dlqProducer != nil {
		c.dlqProducer.Close()
	}
	log.Println("[Kafka] Consumer stopped")
}
