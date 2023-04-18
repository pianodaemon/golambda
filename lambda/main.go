package main

import (
	"context"
	"fmt"
	"os"
	"sync"

	"immortalcrab.com/eventrouter/internal/forwarders"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Buffer size of the channel to take over of the concurrency limit
// by using a channel that can be treated as a kind of pseudo semaphore.
// The number of simultaneous concurrent actions that can be efficiently
// processed will vary depending on the Lambda runtime resources
// https://docs.aws.amazon.com/lambda/latest/dg/gettingstarted-limits.html
const MaxActions = 100

type (
	MessageHandler func(msg *events.SQSMessage) error
)

func mapForEachMessage(records []events.SQSMessage, actOn MessageHandler) (events.SQSEventResponse, error) {

	var wg sync.WaitGroup
	pseudoSem := make(chan struct{}, MaxActions)
	errMsgIdChannel := make(chan string, len(records))
	sqsBatchResponse := events.SQSEventResponse{
		BatchItemFailures: make([]events.SQSBatchItemFailure, 0),
	}

	action := func(msg *events.SQSMessage) {
		pseudoSem <- struct{}{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() {
				<-pseudoSem
			}()

			err := actOn(msg)
			if err != nil {
				emsg := fmt.Sprintf("Exception handling message with id %s", msg.MessageId)
				fmt.Fprintf(os.Stderr, "%s: %s\n", emsg, err.Error())
				errMsgIdChannel <- msg.MessageId
			}
		}()
	}

	for _, msg := range records {
		action(&msg)
	}

	go func() {
		defer close(pseudoSem)
		defer close(errMsgIdChannel)
		wg.Wait()
	}()

	for msgId := range errMsgIdChannel {
		failure := &events.SQSBatchItemFailure{
			ItemIdentifier: msgId,
		}
		sqsBatchResponse.BatchItemFailures = append(sqsBatchResponse.BatchItemFailures, *failure)
	}

	return sqsBatchResponse, nil
}

func handleMessage(msg *events.SQSMessage) (merr error) {

	defer func() {
		if r := recover(); r != nil {
			merr = r.(error)
		}
	}()

	target := forwarders.TargetsLookUp[forwarders.FORWARD_KAFKA]
	target.Forward(msg.Body)

	fmt.Printf("The message %s for event source %s = %s \n", msg.MessageId, msg.EventSource, msg.Body)
	return nil
}

func main() {
	lambda.Start(func(ctx context.Context, sqsEvent events.SQSEvent) (events.SQSEventResponse, error) {

		return mapForEachMessage(sqsEvent.Records, handleMessage)
	})
}
