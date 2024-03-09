package cipher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"io"
	"os"
)

/*
	func ApplePrivateKey() any {
		privateKey, err := readFile(os.Getenv("APPLE_PRIVATE_KEY"))
		if err != nil {
			panic(err)
		}
		block, _ := pem.Decode(privateKey)
		parseResult, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			panic(err)
		}
		return parseResult
	}
*/
func PrivateKey(privateKeyPath string) (key *rsa.PrivateKey) {
	privateKey, err := os.ReadFile(privateKeyPath)
	if err != nil {
		panic(err)
	}
	block, _ := pem.Decode(privateKey)
	key, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	return
}

func PublicKey(publicKeyPath string) (key *rsa.PublicKey) {
	publicKey, err := os.ReadFile(publicKeyPath)
	if err != nil {
		panic(err)
	}
	block, _ := pem.Decode(publicKey)
	parseResult, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	key = parseResult.(*rsa.PublicKey)
	return
}

func EncryptWithKey(value []byte, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Encrypt and prepend nonce
	cipherText := gcm.Seal(nonce, nonce, value, nil)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func DecryptWithKey(cipherText string, key string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(data) < gcm.NonceSize() {
		return nil, errors.New("cipherText too short")
	}

	nonce, cipherTextData := data[:gcm.NonceSize()], data[gcm.NonceSize():]
	plainText, err := gcm.Open(nil, nonce, cipherTextData, nil)
	if err != nil {
		return nil, err
	}

	return plainText, nil
}
