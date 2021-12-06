package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/hill-daniel/drizzle-webhook"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const (
	gitHubSignatureHeader = "X-Hub-Signature-256"
	webhookSecretsPrefix  = "WEBHOOK"
)

// LambdaHandler implements an AWS Lambda handler for an incoming APIGatewayProxyRequest.
type LambdaHandler struct {
	validatePayload drizzle.ValidatePayload
	secretRetriever drizzle.SecretRetriever
	messageQueue    drizzle.Publisher
}

// NewLambdaHandler creates a new LambdaHandler.
func NewLambdaHandler(validator drizzle.ValidatePayload, secretRetriever drizzle.SecretRetriever, messageQueue drizzle.Publisher) *LambdaHandler {
	return &LambdaHandler{
		validatePayload: validator,
		secretRetriever: secretRetriever,
		messageQueue:    messageQueue,
	}
}

// Handle handles an incoming APIGatewayProxyRequest. Tries to validate body with hash signature.
func (lh *LambdaHandler) Handle(ctx context.Context, proxyRequest *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	body := []byte(proxyRequest.Body)
	request := &drizzle.GitHubHookRequest{}
	if err := json.Unmarshal(body, request); err != nil {
		log.Println(errors.Wrapf(err, "failed to parse GitHubHookRequest"))
		return createResponse(http.StatusInternalServerError), nil
	}

	secretValue, err := retrieveWebhookSecret(request.Repository.ID, lh.secretRetriever)
	if err != nil {
		log.Println(err)
		return createResponse(http.StatusInternalServerError), nil
	}

	signature := proxyRequest.Headers[gitHubSignatureHeader]
	if err = lh.validatePayload(secretValue, signature, body); err != nil {
		log.Println(err)
		return createResponse(http.StatusForbidden), nil
	}

	if strings.HasPrefix(request.HeadCommit.Message, "drizzle-pipeline") {
		log.Printf("ignoring commit of drizzle-pipeline with tree id %q", request.HeadCommit.TreeID)
		return createResponse(http.StatusOK), nil
	}

	repository := drizzle.Repository{
		ID:        strconv.FormatInt(request.Repository.ID, 10),
		BranchRef: request.Ref,
		Name:      request.Repository.Name,
		FullName:  request.Repository.FullName,
		Private:   request.Repository.Private,
		URL:       request.Repository.URL,
		CloneURL:  request.Repository.CloneURL,
		Provider:  determineProvider(proxyRequest),
	}
	repositoryMessage, err := json.Marshal(repository)
	if err != nil {
		log.Println(err)
		return createResponse(http.StatusInternalServerError), nil
	}

	if err = lh.messageQueue.Publish(string(repositoryMessage)); err != nil {
		log.Println(err)
		return createResponse(http.StatusInternalServerError), nil
	}

	return createResponse(http.StatusOK), nil
}

func retrieveWebhookSecret(repositoryID int64, secretRetriever drizzle.SecretRetriever) (string, error) {
	secret := fmt.Sprintf("%s_%s_%d", drizzle.GitHubSecretsPrefix, webhookSecretsPrefix, repositoryID)
	secretValue, err := secretRetriever.RetrieveSecret(secret)
	if err != nil {
		return "", errors.Wrapf(err, "failed to retrieve secretValue for %q", secret)
	}
	return secretValue, nil
}

func determineProvider(proxyRequest *events.APIGatewayProxyRequest) string {
	// for now, just GitHub is supported
	return "GitHub"
}

func createResponse(httpStatus int) *events.APIGatewayProxyResponse {
	return &events.APIGatewayProxyResponse{
		StatusCode:        httpStatus,
		Headers:           nil,
		MultiValueHeaders: nil,
		Body:              "",
		IsBase64Encoded:   false,
	}
}
