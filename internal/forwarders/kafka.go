package forwarders

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type (
	DistEventStore struct {
		name        string
		targetTopic string
		config      *kafka.ConfigMap
	}
)

func NewDistEventStore(config *kafka.ConfigMap, targetTopic string) *DistEventStore {

	return &DistEventStore{
		name:        "Kafka confluent",
		targetTopic: targetTopic,
		config:      config,
	}
}

func (self *DistEventStore) GetName() string {
	return self.name
}

func (self *DistEventStore) Forward(payload string) {

	p, err := kafka.NewProducer(self.config)
	if err != nil {
		panic(err)
	}

	defer p.Close()

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	// Produce messages to topic (asynchronously)
	for _, word := range []string{"Welcome", "to", "the", "Confluent", "Kafka", "Golang", "client"} {
		p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &self.targetTopic, Partition: kafka.PartitionAny},
			Value:          []byte(word),
		}, nil)
	}

	// Wait for message deliveries before shutting down
	p.Flush(5 * 1000)
}
