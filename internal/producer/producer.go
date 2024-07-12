package producer

import "github.com/confluentinc/confluent-kafka-go/kafka"

type Producer interface {
	Produce()
	Stop()
}

type SimpleProducer struct {
	producer *kafka.Producer
}

func NewProducer(cfg *kafka.ConfigMap) *SimpleProducer {
	prod, err := kafka.NewProducer(cfg)
	if err != nil {
		panic(err)
	}
	return &SimpleProducer{producer: prod}
}

// Produce dispatches a new kafka message to the
// configured topic
func (s *SimpleProducer) Produce(message *kafka.Message) {

}

// Stop closes the producer
func (s *SimpleProducer) Stop() {
	s.producer.Close()
}
