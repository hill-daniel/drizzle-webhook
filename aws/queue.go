package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/pkg/errors"
)

const (
	messageDelayInSeconds = 5
)

// MessageQueue represents a SQS implementation of an message queue.
type MessageQueue struct {
	queue *sqs.SQS
	url   string
}

// NewMessageQueue creates a new message queue.
func NewMessageQueue(queue *sqs.SQS, url string) *MessageQueue {
	return &MessageQueue{
		queue: queue,
		url:   url,
	}
}

// Publish publishes given message to the queue.
func (mq *MessageQueue) Publish(message string) error {
	if _, err := mq.queue.SendMessage(&sqs.SendMessageInput{
		DelaySeconds: aws.Int64(messageDelayInSeconds),
		MessageBody:  aws.String(message),
		QueueUrl:     aws.String(mq.url),
	}); err != nil {
		return errors.Wrap(err, "failed to publish message to queue")
	}
	return nil
}
