package forwarders

const (
	TARGET_KAFKA = iota
	TARGET_SQS
	TARGET_MAX
)

type (
	Target struct {
		Name    string
		Forward func(payload string)
	}
)

var TargetsLookUp = [TARGET_MAX]*Target{
	&Target{"Kafka confluent", toKafka},
	&Target{"SQS aws", toSqs},
}
