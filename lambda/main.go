package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type (
	MessageHandler func(msg *events.SQSMessage) error
)

func mapForEachMessage(records []events.SQSMessage, actOn MessageHandler) (events.SQSEventResponse, error) {

	sqsBatchResponse := events.SQSEventResponse{
		BatchItemFailures: make([]events.SQSBatchItemFailure, 0),
	}

	for _, msg := range records {
		err := actOn(&msg)
		if err != nil {
			emsg := fmt.Sprintf("Exception handling message with id %s", msg.MessageId)
			fmt.Fprintf(os.Stderr, "%s\n%s\n", emsg, err.Error())
			failure := &events.SQSBatchItemFailure{
				ItemIdentifier: msg.MessageId,
			}
			sqsBatchResponse.BatchItemFailures = append(sqsBatchResponse.BatchItemFailures, *failure)
		}
	}

	return sqsBatchResponse, nil
}

func handleMessage(msg *events.SQSMessage) (merr error) {

	defer func() {
		if r := recover(); r != nil {
			merr = r.(error)
		}
	}()

	fmt.Printf("The message %s for event source %s = %s \n", msg.MessageId, msg.EventSource, msg.Body)
	return nil
}

func main() {
	lambda.Start(func(ctx context.Context, sqsEvent events.SQSEvent) (events.SQSEventResponse, error) {

		return mapForEachMessage(sqsEvent.Records, handleMessage)
	})
}
