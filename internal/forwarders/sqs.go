package forwarders

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type (
	CloudQueue struct {
		name     string
		client   *sqs.Client
		queue    string
		queueUrl *string
	}
)

func NewCloudQueue(queue string, cfg aws.Config) *CloudQueue {

	var client *sqs.Client = sqs.NewFromConfig(cfg)
	gQInput := &sqs.GetQueueUrlInput{
		QueueName: &queue,
	}

	ctx := context.TODO()
	result, err := client.GetQueueUrl(ctx, gQInput)
	if err != nil {
		panic(err)
	}

	return &CloudQueue{
		name:     "SQS aws",
		client:   client,
		queue:    queue,
		queueUrl: result.QueueUrl,
	}
}

func (self *CloudQueue) GetName() string {
	return self.name
}

func (self *CloudQueue) Forward(payload string) {

	sMInput := &sqs.SendMessageInput{
		DelaySeconds: 10,
		MessageAttributes: map[string]types.MessageAttributeValue{
			"Title": {
				DataType:    aws.String("String"),
				StringValue: aws.String("The Whistler"),
			},
			"Author": {
				DataType:    aws.String("String"),
				StringValue: aws.String("John Grisham"),
			},
			"WeeksOn": {
				DataType:    aws.String("Number"),
				StringValue: aws.String("6"),
			},
		},
		MessageBody: aws.String("Information about the NY Times fiction bestseller for the week of 12/11/2016."),
		QueueUrl:    self.queueUrl,
	}

	ctx := context.TODO()
	resp, err := self.client.SendMessage(ctx, sMInput)
	if err != nil {
		fmt.Println("Got an error sending the message:")
		fmt.Println(err)
		return
	}

	fmt.Println("Sent message with ID: " + *resp.MessageId)
}
