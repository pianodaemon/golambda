package forwarders

const (
	FORWARD_KAFKA = iota
	FORWARD_TARGET_MAX
)

type (
	Target struct {
		Name    string
		Forward func(payload string)
	}
)

var TargetsLookUp = [FORWARD_TARGET_MAX]*Target{
	&Target{"Kafka concluent", toKafka},
}
