package aws_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/hill-daniel/drizzle-webhook"
	"github.com/hill-daniel/drizzle-webhook/aws"
	"github.com/hill-daniel/drizzle-webhook/mocks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

func TestLambdaHandler_Handle_should_validate_signature_with_payload(t *testing.T) {
	var validator drizzle.ValidatePayload
	validator = func(secret, signature string, payload []byte) error {
		assert.Equal(t, "super-secret", secret)
		assert.Equal(t, "sha256=signature", signature)
		assert.Equal(t, "{}", string(payload))
		return nil
	}
	secretsManager := &mocks.SecretsManager{GetSecretFunc: func(secretID string) (string, error) {
		return "super-secret", nil
	}}
	messageQueue := &mocks.MessageQueue{PublishFunc: func(message string) error {
		return nil
	}}
	handler := aws.NewLambdaHandler(validator, secretsManager, messageQueue)
	request := createProxyRequest()
	request.Headers["X-Hub-Signature-256"] = "sha256=signature"

	response, err := handler.Handle(context.TODO(), request)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func TestLambdaHandler_Handle_should_return_error_on_validation_failure(t *testing.T) {
	var failingValidator drizzle.ValidatePayload
	failingValidator = func(secret, signature string, payload []byte) error {
		return fmt.Errorf("failed to validate")
	}
	secretsManager := &mocks.SecretsManager{GetSecretFunc: func(secretID string) (string, error) {
		return "", nil
	}}
	messageQueue := &mocks.MessageQueue{PublishFunc: func(message string) error {
		return nil
	}}
	handler := aws.NewLambdaHandler(failingValidator, secretsManager, messageQueue)
	request := createProxyRequest()

	response, err := handler.Handle(context.TODO(), request)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, response.StatusCode)
}

func Test_should_ignore_request_with_drizzle_head_commit(t *testing.T) {
	request, err := createRequest(t)
	request.HeadCommit.Message = "drizzle-pipeline fix some stuff"
	body, err := json.Marshal(request)
	assert.NoError(t, err)
	var invalidatingRequestValidator drizzle.ValidatePayload
	invalidatingRequestValidator = func(secret, signature string, payload []byte) error {
		return nil
	}
	secretsManager := &mocks.SecretsManager{GetSecretFunc: func(secretID string) (string, error) {
		return "", nil
	}}
	messageQueue := &mocks.MessageQueue{PublishFunc: func(message string) error {
		return errors.New("message queue should not have been called")
	}}
	handler := aws.NewLambdaHandler(invalidatingRequestValidator, secretsManager, messageQueue)
	proxyRequest := createProxyRequest()
	proxyRequest.Body = string(body)

	response, err := handler.Handle(context.TODO(), proxyRequest)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func createRequest(t *testing.T) (*drizzle.GitHubHookRequest, error) {
	f, err := os.Open("../testdata/hook/request.json")
	assert.NoError(t, err)
	fileContent, err := ioutil.ReadAll(f)
	assert.NoError(t, err)
	request := &drizzle.GitHubHookRequest{}
	err = json.Unmarshal(fileContent, request)
	assert.NoError(t, err)
	return request, err
}

func createProxyRequest() *events.APIGatewayProxyRequest {
	request := &events.APIGatewayProxyRequest{}
	request.Headers = make(map[string]string)
	request.Body = "{}"
	return request
}
