package forwarders

import (
	"context"

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

	TargetsLookUp[TARGET_KAFKA] = NewDistEventStore(&kafka.ConfigMap{
		"bootstrap.servers":            "localhost",
		"queue.buffering.max.messages": "1",
		"queue.buffering.max.ms":       "1",
	})

	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err)
	}

	TargetsLookUp[TARGET_SQS] = NewCloudQueue("polito", cfg)
}
