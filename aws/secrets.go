package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/pkg/errors"
)

// SecretsManager represents an AWS SecretsManager.
type SecretsManager struct {
	sm *secretsmanager.SecretsManager
}

// NewSecretManager creates a new SecretsManager.
func NewSecretManager(secretsManager *secretsmanager.SecretsManager) *SecretsManager {
	return &SecretsManager{
		sm: secretsManager,
	}
}

// RetrieveSecret retrieves a secret value from SecretsManager.
func (s *SecretsManager) RetrieveSecret(secretID string) (string, error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretID),
	}
	value, err := s.sm.GetSecretValue(input)
	if err != nil {
		return "", errors.Wrapf(err, "failed to retrieve value of secret %q", secretID)
	}
	return *value.SecretString, nil
}
