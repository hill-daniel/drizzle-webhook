package drizzle

const (
	// GitHubSecretsPrefix used to retrieve GitHub secrets.
	GitHubSecretsPrefix = "GITHUB"
)

// ValidatePayload validates a given payload with signature and secret.
type ValidatePayload func(secret, signature string, payload []byte) error

//SecretRetriever retrieves secrets.
type SecretRetriever interface {
	// RetrieveSecret retrieves a secret for given ID.
	RetrieveSecret(secretID string) (string, error)
}

// Publisher is an interface for publishing messages.
type Publisher interface {
	Publish(message string) error
}

// Repository represents a git repository.
type Repository struct {
	ID        string `json:"id"`
	BranchRef string `json:"branch_ref"`
	Name      string `json:"name"`
	FullName  string `json:"full_name"`
	Private   bool   `json:"private"`
	URL       string `json:"url"`
	CloneURL  string `json:"clone_url"`
	Provider  string `json:"provider"`
}
