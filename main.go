package main

import (
	"context"
	"fmt"
	"time"

	"github.com/tandy9527/js-util/kafka"
)

func main() {
	brokers := []string{"localhost:9092"}
	topics := []string{"test-topic"}
	groupID := "group1"
	dlqTopic := "test-topic-dlq"

	// 1️⃣ 初始化 Producer
	kafka.InitProducer(brokers, topics[0]) // 第一个 topic 用作 Producer 默认 topic

	// 2️⃣ 初始化 Consumer
	handler := func(ctx context.Context, key, value []byte) error {
		fmt.Printf("[Consumer] key=%s, value=%s\n", string(key), string(value))
		// 模拟随机失败
		if time.Now().UnixNano()%2 == 0 {
			return fmt.Errorf("simulated handler error")
		}
		return nil
	}

	kafka.InitConsumer(
		kafka.Config{
			Brokers:        brokers,
			GroupID:        groupID,
			Topics:         topics,
			MinBytes:       10e3,
			MaxBytes:       10e6,
			CommitInterval: 1000,
		},
		handler,
		3,           // retryCount
		time.Second, // retryDelay
		dlqTopic,    // DLQ topic
	)
	select {}

}
