package main

import (
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/symonk/kafka-demo/internal/consumer"
	"github.com/symonk/kafka-demo/internal/models"
	"github.com/symonk/kafka-demo/internal/producer"
)

func main() {
	topic := "heroes"
	prodCfg := kafka.ConfigMap{"bootstrap.servers": "localhost"}
	consumerCfg := kafka.ConfigMap{"bootstrap.servers": "localhost", "group.id": "demo"}
	producer := producer.New(&prodCfg)
	defer producer.Stop()
	consumer := consumer.New(&consumerCfg)
	defer consumer.Stop()

	for i := range 10000 {
		value := models.Superhero{Id: i, Name: fmt.Sprintf("name%d", i), Power: i * 100, Melee: false}
		bs, err := json.Marshal(value)
		if err != nil {
			panic(err)
		}

		err = producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          bs,
			Headers:        []kafka.Header{{Key: "demo-header", Value: []byte("demo")}},
		})
		if err != nil {
			fmt.Println(err)
		}
	}
	producer.Wait()
}
