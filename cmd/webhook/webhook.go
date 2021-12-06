package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/hill-daniel/drizzle-webhook/aws"
	"github.com/hill-daniel/drizzle-webhook/crypto"
	"github.com/pkg/errors"
	"log"
	"os"
)

const (
	messageQueueURL = "MESSAGE_QUEUE_URL"
)

func main() {
	awsSession, err := session.NewSession()
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to create aws session"))
	}
	secretsManager := aws.NewSecretManager(secretsmanager.New(awsSession))
	queueURL := os.Getenv(messageQueueURL)
	sqsQ := sqs.New(awsSession)
	queue := aws.NewMessageQueue(sqsQ, queueURL)
	handler := aws.NewLambdaHandler(crypto.ValidatePayload, secretsManager, queue)
	lambda.Start(handler.Handle)
}
