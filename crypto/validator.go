package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/pkg/errors"
	"strings"
)

const (
	sha256Prefix       = "sha256"
	signatureSeparator = "="
	numberOfParts      = 2
)

// ValidatePayload checks if the payload's calculated hash is equal to given signature.
func ValidatePayload(secret, signature string, payload []byte) error {
	sigParts := strings.SplitN(signature, signatureSeparator, numberOfParts)
	if len(sigParts) != numberOfParts || sigParts[0] != sha256Prefix {
		return errors.New("failed to parse signature")
	}

	decoded, err := hex.DecodeString(sigParts[1])
	if err != nil {
		return errors.Wrap(err, "failed to decode signature")
	}

	calculated := calcHash(secret, payload)

	if !hmac.Equal(calculated, decoded) {
		return errors.New("failed to validate. Calculated hash does not match signature")
	}
	return nil
}

// calcHash computes the hash of payload's body according to the webhook's secret token
// see https://developer.github.com/webhooks/securing/#validating-payloads-from-github
func calcHash(secret string, payload []byte) []byte {
	hm := hmac.New(sha256.New, []byte(secret))
	hm.Write(payload)
	return hm.Sum(nil)
}
