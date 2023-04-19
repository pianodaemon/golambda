package forwarders

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

const (
	TARGET_KAFKA = iota
	TARGET_SQS
	TARGET_MAX
)

type (
	Target interface {
		GetName() string
		Forward(payload string)
	}
)

var TargetsLookUp = make([]Target, TARGET_MAX)

func init() {

	getEnv := func(key, fallback string) string {
		if value, ok := os.LookupEnv(key); ok {
			return value
		}
		return fallback
	}

	kafkaTopic := getEnv("KAFKA_TARGET_TOPIC", "internal_test")
	kafkaServers := getEnv("KAFKA_SERVERS", "localhost:9092")
	configMap := &kafka.ConfigMap{
		"bootstrap.servers":            kafkaServers,
		"queue.buffering.max.messages": "1",
		"queue.buffering.max.ms":       "1",
	}
	TargetsLookUp[TARGET_KAFKA] = NewTargetKafka(configMap, kafkaTopic)

	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err)
	}

	TargetsLookUp[TARGET_SQS] = NewTargetSQS("polito", cfg)
}
