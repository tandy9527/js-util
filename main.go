package main

import (
	"context"
	"fmt"
	"js-util/kafka"
	"time"
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
	// consumer := kafka.GetConsumer()
	// producer := kafka.GetProducer()

	// // 3️⃣ Producer 循环发送消息
	// go func() {
	// 	i := 0
	// 	for {
	// 		key := []byte(fmt.Sprintf("key-%d", i%3))
	// 		value := []byte(fmt.Sprintf("msg-%d %s", i, time.Now().Format(time.RFC3339)))
	// 		if err := producer.SendMessage(key, value); err != nil {
	// 			log.Printf("[Producer] send failed: %v", err)
	// 		} else {
	// 			log.Printf("[Producer] sent: %s", string(value))
	// 		}
	// 		i++
	// 		time.Sleep(time.Second)
	// 	}
	// }()

	// // 4️⃣ 等待退出信号
	// sigs := make(chan os.Signal, 1)
	// signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	// <-sigs

	// log.Println("Shutting down...")
	// consumer.Stop()
	// producer.Close()
	// log.Println("Exit")
	// 3️⃣ 初始化 ProducerRunner
	// kafka.InitConsumer(time.Second, func(i int) ([]byte, []byte) {
	// 	key := []byte(fmt.Sprintf("key-%d", i%3))
	// 	value := []byte(fmt.Sprintf("msg-%d %s", i, time.Now().Format(time.RFC3339)))
	// 	return key, value
	// })

	// runner := kafka.GetConsumer()
	// runner.Start()

	// // 阻塞等待退出
	// runner.WaitForSignal()
}
