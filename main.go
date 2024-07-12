package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/symonk/kafka-demo/internal/consumer"
	"github.com/symonk/kafka-demo/internal/models"
	"github.com/symonk/kafka-demo/internal/producer"
)

func main() {
	topic := "heroes"
	prodCfg := kafka.ConfigMap{"bootstrap.servers": "localhost"}
	consumerCfg := kafka.ConfigMap{"bootstrap.servers": "localhost", "group.id": "demo", "statistics.interval.ms": 5000, "auto.offset.reset": "earliest", "enable.auto.commit": false}
	producer := producer.New(&prodCfg)
	defer producer.Stop()
	consumer := consumer.New(&consumerCfg)
	defer consumer.Stop()

	for i := range 1_000_000 {
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

	var wg sync.WaitGroup
	wg.Add(1_000_000)

	// start a group of consumers, we have ~10 partitions configured so, 10 will do!
	for range 10 {
		go func() {
			c, err := kafka.NewConsumer(&consumerCfg)
			if err != nil {
				panic(err)
			}
			// no callback for rebalancing for now.
			err = c.SubscribeTopics([]string{topic}, nil)
			if err != nil {
				panic(err)
			}

			// each consumer should get around 1/10th messages with a null key
			for range 100_000 {
				msg, err := c.ReadMessage(time.Second)
				if err == nil {
					fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
					_, err := c.Commit()
					if err != nil {
						fmt.Printf("error %s", err)
					}
				} else if !err.(kafka.Error).IsTimeout() {
					fmt.Printf("Consumer error: %v (%v)\n", err, msg)
				}
				wg.Done()
			}
		}()
	}

	producer.Wait()
	wg.Wait()
}
