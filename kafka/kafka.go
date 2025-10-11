package kafka

// import (
// 	"github.com/segmentio/kafka-go"
// )

// type KafkaProducer struct {
// 	writer     *kafka.Writer
// 	retryCount int
// 	batchSize  int
// }

// type KafkaConsumer struct {
// 	readers []*kafka.Reader
// 	cfg     Config
// }

// // 初始化生产者
// func NewProducer(cfg Config) *KafkaProducer {
// 	writer := kafka.NewWriter(kafka.WriterConfig{
// 		Brokers:  cfg.Brokers,
// 		Topic:    cfg.Topic,
// 		Balancer: &kafka.LeastBytes{},
// 	})
// 	return &KafkaProducer{
// 		writer:     writer,
// 		retryCount: cfg.RetryCount,
// 		batchSize:  cfg.BatchSize,
// 	}
// }

// // 初始化消费者
// func NewConsumer(cfg Config) *KafkaConsumer {
// 	readers := make([]*kafka.Reader, cfg.ConsumerCount)
// 	for i := 0; i < cfg.ConsumerCount; i++ {
// 		readers[i] = kafka.NewReader(kafka.ReaderConfig{
// 			Brokers: cfg.Brokers,
// 			Topic:   cfg.Topic,
// 			GroupID: cfg.GroupID,
// 		})
// 	}
// 	return &KafkaConsumer{readers: readers, cfg: cfg}
// }

// // 关闭生产者
// func (p *KafkaProducer) Close() {
// 	p.writer.Close()
// }

// // 关闭消费者
// func (c *KafkaConsumer) Close() {
// 	for _, r := range c.readers {
// 		r.Close()
// 	}
// }
