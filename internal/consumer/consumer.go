package consumer

import "github.com/confluentinc/confluent-kafka-go/v2/kafka"

type Consumer interface {
	Consume()
	Stop()
}

type SimpleConsumer struct {
	consumer *kafka.Consumer
}

func New(cfg *kafka.ConfigMap) *SimpleConsumer {
	c, err := kafka.NewConsumer(cfg)
	if err != nil {
		panic(err)
	}

	return &SimpleConsumer{consumer: c}
}

// Stop closes the consumer
func (s *SimpleConsumer) Stop() {
	// TODO: err handling
	s.consumer.Close()
}
