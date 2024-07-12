package producer

import (
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Producer interface {
	Produce()
	Stop()
}

type SimpleProducer struct {
	producer *kafka.Producer
}

func New(cfg *kafka.ConfigMap) *SimpleProducer {
	prod, err := kafka.NewProducer(cfg)
	if err != nil {
		panic(err)
	}
	s := &SimpleProducer{producer: prod}
	go s.Monitor()
	return s
}

// Produce dispatches a new kafka message to the
// configured topic
func (s *SimpleProducer) Produce(message *kafka.Message) error {
	err := s.producer.Produce(message, nil)
	if err != nil {
		if err.(kafka.Error).Code() == kafka.ErrQueueFull {
			time.Sleep(time.Second)
		}
		fmt.Printf("Failed to produce message: %v\n", err)
	}
	return err
}

func (s *SimpleProducer) Wait() {
	for s.producer.Flush(10_000) > 0 {
		fmt.Print("still waiting to flush all messages from the producer")

	}
}

// Monitor listens to the producers
func (s *SimpleProducer) Monitor() {
	for event := range s.producer.Events() {
		switch eventType := event.(type) {
		case *kafka.Message:
			// the message delivery report, indicating success or
			// permanent failure after retries have been exhausted.
			message := eventType
			if message.TopicPartition.Error != nil {
				fmt.Printf("Delivery failed: %v\n", message.TopicPartition.Error)
			}
			fmt.Printf("Delivered message to topic %s, partition: [%d] at offset: %v\n", *message.TopicPartition.Topic, message.TopicPartition.Partition, message.TopicPartition.Offset)
		case kafka.Error:
			// client instance-level errors
			// these are auto retried and should be considered informational
			fmt.Printf("Error: %v\n", eventType)
		default:
			fmt.Printf("Ignored event: %s\n", eventType)
		}
	}

}

// Stop closes the producer
func (s *SimpleProducer) Stop() {
	s.producer.Close()
}
