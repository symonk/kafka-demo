package main

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/symonk/kafka-demo/internal/consumer"
	"github.com/symonk/kafka-demo/internal/producer"
)

func main() {
	fmt.Println("in main")
	cfg := kafka.ConfigMap{"bootstrap.servers": "localhost"}
	prod := producer.New(&cfg)
	defer prod.Stop()
	consumer := consumer.New(&cfg)
	defer consumer.Stop()
}
