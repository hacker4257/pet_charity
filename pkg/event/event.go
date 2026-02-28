package event

import (
	"context"
	"encoding/json"
	"time"

	"github.com/hacker4257/pet_charity/pkg/logger"
	"github.com/segmentio/kafka-go"
)

type Event struct {
	Action    string `json:"action"`
	UserID    uint   `json:"user_id"`
	RelatedID uint   `json:"related_id"`
}

type Handler func(e Event)

var (
	writer   *kafka.Writer
	reader   *kafka.Reader
	handlers []Handler
	topic    string
)

func Init(brokers []string, topicName string) {
	topic = topicName
	writer = &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 10 * time.Millisecond,
	}
}

// Subscribe 注册监听
func Subscribe(h Handler) {
	handlers = append(handlers, h)
}

// Publish 发布事件到 Kafka
func Publish(action string, userID uint, relatedID uint) {
	e := Event{Action: action, UserID: userID, RelatedID: relatedID}
	data, err := json.Marshal(e)
	if err != nil {
		logger.Errorf("[event] marshal failed: %v", err)
		return
	}
	go func() {
		err := writer.WriteMessages(context.Background(), kafka.Message{
			Value: data,
		})
		if err != nil {
			logger.Warnf("kafka publish failed: %v", err)
		}
	}()
}

// StartConsumer 启动 Kafka 消费者（main 里调用）

func StartConsumer(brokers []string, topicName, groupID string) {
	reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topicName,
		GroupID:  groupID,
		MinBytes: 1,
		MaxBytes: 10e6,
	})

	go func() {
		for {
			msg, err := reader.ReadMessage(context.Background())
			if err != nil {
				logger.Warnf("[event] kafka read failed: %v", err)
				time.Sleep(time.Second) // backoff
				continue
			}
			var e Event
			if err := json.Unmarshal(msg.Value, &e); err != nil {
				logger.Warnf("[event] unmarshal failed: %v", err)
				continue
			}
			//分发给所有监听者
			for _, h := range handlers {
				func() {
					defer func() {
						if r := recover(); r != nil {
							logger.Errorf("[event] handler panic: %v", r)
						}
					}()
					h(e)
				}()
			}
		}
	}()
}

func Close() {
	if writer != nil {
		writer.Close()
	}
	if reader != nil {
		reader.Close()
	}
}
