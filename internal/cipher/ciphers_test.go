package cipher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestPrivateKeyLoading tests loading of a RSA private key from file
func TestPrivateKeyLoading(t *testing.T) {
	privateKeyPath := "../../vault/example/jwt.private"
	key := PrivateKey(privateKeyPath)
	assert.NotNil(t, key, "PrivateKey should not be nil")
}

// TestPublicKeyLoading tests loading of a RSA public key from file
func TestPublicKeyLoading(t *testing.T) {
	publicKeyPath := "../../vault/example/jwt.public"
	key := PublicKey(publicKeyPath)
	assert.NotNil(t, key, "PublicKey should not be nil")
}

// TestEncryptDecryptWithKey tests the encryption and decryption flow
func TestEncryptDecryptWithKey(t *testing.T) {
	value := "hello world"
	key := "abcd"
	encryptedValue, err := EncryptWithKey([]byte(value), key)
	assert.NoError(t, err, "EncryptWithKey should not error")
	assert.NotEmpty(t, encryptedValue, "encryptedValue should not be empty")

	decryptedValue, err := DecryptWithKey(encryptedValue, key)
	assert.NoError(t, err, "DecryptWithKey should not error")
	assert.Equal(t, value, string(decryptedValue), "decryptedValue should match original value")
}
