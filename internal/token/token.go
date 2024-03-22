package token

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/cristalhq/jwt/v4"
	"github.com/google/uuid"
	"github.com/usercoredev/usercore/internal/cipher"
	"time"
)

type Token jwt.RegisteredClaims

type claimsKey string

var Claims claimsKey = "claims"

type DefaultToken struct {
	AccessToken  string
	RefreshToken string
}

type PublicPrivateKey struct {
	PublicKey  *rsa.PublicKey
	PrivateKey *rsa.PrivateKey
}

type Settings struct {
	Scheme             string
	Issuer             string
	Audience           string
	PrivateKeyPath     string
	PublicKeyPath      string
	RefreshTokenExpire time.Duration
	AccessTokenExpire  time.Duration
	PublicPrivateKey   PublicPrivateKey
	Verifier           jwt.Verifier
	Signer             jwt.Signer
}

var options *Settings

func (s *Settings) Setup() {
	publicKey := cipher.PublicKey(s.PublicKeyPath)
	privateKey := cipher.PrivateKey(s.PrivateKeyPath)
	verifier, err := jwt.NewVerifierPS(jwt.PS512, publicKey)
	if err != nil {
		panic(err)
	}
	signer, err := jwt.NewSignerPS(jwt.PS512, privateKey)
	if err != nil {
		panic(err)
	}
	s.Verifier = verifier
	s.Signer = signer
	options = s
}

func CreateRefreshToken(uuid uuid.UUID) (string, *time.Time) {
	buffer := make([]byte, 16)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", nil
	}

	random := hex.EncodeToString(buffer)

	if err != nil {
		panic(err)
	}
	refreshTokenExpireTime := time.Now().Add(options.RefreshTokenExpire)
	var content = uuid.String() + fmt.Sprint(time.Now().Unix()) + random
	hashed := sha256.New()
	hashed.Write([]byte(content))
	refreshTokenString := hashed.Sum(nil)
	return hex.EncodeToString(refreshTokenString), &refreshTokenExpireTime
}

func CreateJWT(userId uuid.UUID) (string, error) {
	registeredClaims := &jwt.RegisteredClaims{
		Issuer:    options.Issuer,
		ID:        userId.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(options.AccessTokenExpire)),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Audience:  jwt.Audience{options.Audience},
	}
	builder := jwt.NewBuilder(options.Signer)
	token, err := builder.Build(registeredClaims)
	if err != nil {
		return "", err
	}
	return token.String(), nil
}
