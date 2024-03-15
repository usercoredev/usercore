package token

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/cristalhq/jwt/v4"
	"github.com/google/uuid"
	"github.com/usercoredev/usercore/app/responses"
	"github.com/usercoredev/usercore/utils/cipher"
	"os"
	"strconv"
	"time"
)

type Token jwt.RegisteredClaims

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
	RefreshTokenExpire string
	AccessTokenExpire  string
	PublicPrivateKey   PublicPrivateKey
}

var options *Settings

func (o *Settings) Setup() {
	publicKey := cipher.PublicKey(o.PublicKeyPath)
	privateKey := cipher.PrivateKey(o.PrivateKeyPath)
	o.PublicPrivateKey = PublicPrivateKey{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}
	options = o
}

func CreateRefreshToken(uuid uuid.UUID) (string, *time.Time) {
	buffer := make([]byte, 16)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", nil
	}

	random := hex.EncodeToString(buffer)

	refreshTokenExpire, err := strconv.Atoi(options.RefreshTokenExpire)
	if err != nil {
		panic(err)
	}
	refreshTokenExpireTime := time.Now().Add(time.Duration(refreshTokenExpire) * time.Minute)
	var content = uuid.String() + fmt.Sprint(time.Now().Unix()) + random
	hashed := sha256.New()
	hashed.Write([]byte(content))
	refreshTokenString := hashed.Sum(nil)
	return hex.EncodeToString(refreshTokenString), &refreshTokenExpireTime
}

func CreateJWT(userId uuid.UUID) (string, error) {
	signer, err := jwt.NewSignerPS(jwt.PS512, options.PublicPrivateKey.PrivateKey)
	if err != nil {
		return "", err
	}

	accessTokenExpire, err := strconv.Atoi(options.AccessTokenExpire)
	if err != nil {
		panic(err)
	}
	claims := &jwt.RegisteredClaims{
		Issuer:    os.Getenv("APP_NAME"),
		ID:        userId.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(accessTokenExpire) * time.Minute)),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Audience:  jwt.Audience{options.Audience},
	}
	builder := jwt.NewBuilder(signer)
	token, err := builder.Build(claims)
	if err != nil {
		return "", err
	}
	return token.String(), nil
}

func VerifyJWT(token string) (*Token, error) {
	verifier, err := jwt.NewVerifierPS(jwt.PS512, options.PublicPrivateKey.PublicKey)
	if err != nil {
		return nil, fmt.Errorf(responses.InvalidToken)
	}

	tokenBytes := []byte(token)
	newToken, parseErr := jwt.Parse(tokenBytes, verifier)
	if parseErr != nil {
		return nil, fmt.Errorf(responses.InvalidToken)
	}

	err = verifier.Verify(newToken)
	if err != nil {
		return nil, fmt.Errorf(responses.InvalidToken)
	}

	var newClaims jwt.RegisteredClaims
	err = json.Unmarshal(newToken.Claims(), &newClaims)
	if err != nil {
		return nil, fmt.Errorf(responses.TokenMalformed)
	}

	// or parse only claims
	err = jwt.ParseClaims(tokenBytes, verifier, &newClaims)
	if err != nil {
		return nil, fmt.Errorf(responses.TokenMalformed)
	}

	// verify claims as you wish
	var isValid = newClaims.IsValidAt(time.Now())
	if !isValid {
		return nil, fmt.Errorf(responses.TokenExpired)
	}
	return (*Token)(&newClaims), nil
}
