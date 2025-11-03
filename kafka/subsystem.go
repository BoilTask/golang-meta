package kafka

import (
	"context"
	"encoding/json"
	"log"
	"meta/engine"
	metaerror "meta/meta-error"
	metapanic "meta/meta-panic"
	metastring "meta/meta-string"
	"meta/subsystem"
	"time"

	"github.com/segmentio/kafka-go"
)

type Subsystem struct {
	subsystem.Subsystem
	GetConfig func() *Config

	Addr string `yaml:"addr"`
}

func GetSubsystem() *Subsystem {
	if thisSubsystem := engine.GetSubsystem[*Subsystem](); thisSubsystem != nil {
		return thisSubsystem.(*Subsystem)
	}
	return nil
}

func (s *Subsystem) GetName() string {
	return "Kafka"
}

func (s *Subsystem) Start() error {
	config := s.GetConfig()
	if config == nil {
		return metaerror.New("kafka config is nil")
	}
	s.Addr = config.Addr
	return nil
}

func (s *Subsystem) Stop() error {
	return nil
}

// ProduceMessage 推送消息
func (s *Subsystem) ProduceMessage(ctx context.Context, topic string, key string, value interface{}) error {
	// 转换对象为 JSON
	valueJSON, err := json.Marshal(value)
	if err != nil {
		log.Printf("Error marshalling value: %v", err)
		return err
	}

	// 每次推送时创建新的 Kafka Writer 实例，避免共享问题
	writer := &kafka.Writer{
		// Kafka 地址，根据实际配置修改
		Addr:     kafka.TCP(s.Addr),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	defer func(writer *kafka.Writer) {
		err := writer.Close()
		if err != nil {
			metapanic.ProcessError(metaerror.Wrap(err))
		}
	}(writer)

	// 创建消息
	message := kafka.Message{
		Key:   []byte(key),
		Value: valueJSON,
	}

	// 推送消息
	if err := writer.WriteMessages(ctx, message); err != nil {
		log.Printf("Failed to produce message: %v", err)
		return err
	}

	log.Printf("Message produced to topic %s with key %s", topic, key)
	return nil
}

func (s *Subsystem) ProduceMessageSimple(ctx context.Context, topic string, value interface{}) error {
	return s.ProduceMessage(ctx, topic, metastring.GetRandomString(5), value)
}

// Subscribe 订阅消息
func (s *Subsystem) Subscribe(groupId string, topic string, callback func(key, value string) error) error {
	for {
		reader := kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{s.Addr},
			GroupID: groupId,
			Topic:   topic,
		})

		for {
			message, err := reader.ReadMessage(context.Background())
			if err != nil {
				metapanic.ProcessError(metaerror.Wrap(err))
				break
			}
			if err := callback(string(message.Key), string(message.Value)); err != nil {
				metapanic.ProcessError(metaerror.Wrap(err, "error handling message, topic: %s", topic))
			}
		}

		_ = reader.Close()

		// 暂时先认为是种错误
		metapanic.ProcessError(metaerror.New("kafka reader closed, topic: %s", topic))
		// 等待一段时间后重新订阅
		time.Sleep(10 * time.Second)
	}
}
