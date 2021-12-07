package crypto_test

import (
	"github.com/hill-daniel/drizzle-webhook/crypto"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

const (
	// don't worry, it's an old secret :P
	testSecret = "28E87E52-8EA4-47F9-BD2B-D55AE352009C"
)

func TestValidatePayload_should_validate_given_message(t *testing.T) {
	f, err := os.Open("../testdata/hook/request.json")
	assert.NoError(t, err)
	body, err := ioutil.ReadAll(f)
	assert.NoError(t, err)

	err = crypto.ValidatePayload(testSecret, "sha256=886683835f03ba1719007744f9111d5f09c3f7505b6333aa4a546af15a2cbc81", body)

	assert.NoError(t, err)
}

func TestValidatePayload_should_fail_on_not_matching_hashes(t *testing.T) {
	f, err := os.Open("../testdata/hook/request.json")
	assert.NoError(t, err)
	body, err := ioutil.ReadAll(f)
	assert.NoError(t, err)

	err = crypto.ValidatePayload(testSecret, "sha256=05b6333aa4a546af15", body)

	assert.EqualErrorf(t, err, "failed to validate. Calculated hash does not match signature", "")
}

func TestValidatePayload_should_fail_on_invalid_signature(t *testing.T) {
	err := crypto.ValidatePayload("someSecret", "invalidSignature", []byte{})

	assert.EqualErrorf(t, err, "failed to parse signature", "")
}

func TestValidatePayload_should_fail_on_empty_signature(t *testing.T) {
	err := crypto.ValidatePayload("someSecret", "", []byte{})
	assert.EqualErrorf(t, err, "failed to parse signature", "")
}

func TestValidatePayload_should_fail_if_no_sha256_prefix(t *testing.T) {
	err := crypto.ValidatePayload("someSecret", "s9j=886683835f03ba1719007744f9111d5f09c3f7505b6333aa4a546af15a2cbc81", []byte{})

	assert.EqualErrorf(t, err, "failed to parse signature", "")
}

func TestValidatePayload_should_fail_if_signature_is_not_hex_encoded(t *testing.T) {
	err := crypto.ValidatePayload("someSecret", "sha256=ZZZ28763427", []byte{})

	assert.EqualErrorf(t, err, "failed to decode signature: encoding/hex: invalid byte: U+005A 'Z'", "")
}
