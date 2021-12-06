package mocks

// SecretsManager mock.
type SecretsManager struct {
	GetSecretFunc func(secretID string) (string, error)
}

// RetrieveSecret GetValue mock.
func (sm *SecretsManager) RetrieveSecret(secretID string) (string, error) {
	return sm.GetSecretFunc(secretID)
}

// MessageQueue mock.
type MessageQueue struct {
	PublishFunc func(message string) error
}

// Publish mock.
func (mq *MessageQueue) Publish(message string) error {
	return mq.PublishFunc(message)
}
