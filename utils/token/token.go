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

type DefaultToken struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
}

type tokenOptions struct {
	Scheme               string
	AccessTokenExpire    time.Time
	AccessTokenExpireIn  int
	RefreshTokenExpire   time.Time
	RefreshTokenExpireIn int
}

type PublicPrivateKey struct {
	PublicKey  *rsa.PublicKey
	PrivateKey *rsa.PrivateKey
}

var publicPrivateKey PublicPrivateKey

func SetPublicPrivateKey(publicKeyPath, privateKeyPath string) {
	publicKey := cipher.PublicKey(publicKeyPath)
	privateKey := cipher.PrivateKey(privateKeyPath)
	publicPrivateKey = PublicPrivateKey{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}
}

var options tokenOptions

func SetOptions(accessTokenExpire, refreshTokenExpire, scheme string) {
	accessTokenExpireInMinute, err := strconv.Atoi(accessTokenExpire)
	if err != nil {
		panic(err)
	}

	refreshTokenExpireInMinute, err := strconv.Atoi(refreshTokenExpire)
	if err != nil {
		panic(err)
	}

	options = tokenOptions{
		Scheme:               scheme,
		AccessTokenExpire:    time.Now().Add(time.Duration(accessTokenExpireInMinute) * time.Second),
		AccessTokenExpireIn:  accessTokenExpireInMinute,
		RefreshTokenExpire:   time.Now().Add(time.Duration(refreshTokenExpireInMinute) * time.Second),
		RefreshTokenExpireIn: refreshTokenExpireInMinute,
	}
}

func CreateRefreshToken(uuid uuid.UUID) (string, *time.Time) {
	buffer := make([]byte, 16)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", nil
	}

	random := hex.EncodeToString(buffer)

	var content = uuid.String() + fmt.Sprint(time.Now().Unix()) + random
	hashed := sha256.New()
	hashed.Write([]byte(content))
	refreshTokenString := hashed.Sum(nil)
	return hex.EncodeToString(refreshTokenString), &options.RefreshTokenExpire
}

func CreateJWT(userId uuid.UUID) (string, int64, error) {
	signer, err := jwt.NewSignerPS(jwt.PS512, publicPrivateKey.PrivateKey)
	if err != nil {
		return "", 0, err
	}

	claims := &jwt.RegisteredClaims{
		Issuer:    os.Getenv("COMPANY_NAME"),
		Subject:   os.Getenv("APP_NAME"),
		ID:        userId.String(),
		ExpiresAt: jwt.NewNumericDate(options.AccessTokenExpire),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Audience:  []string{"admin"},
	}
	builder := jwt.NewBuilder(signer)
	token, err := builder.Build(claims)
	if err != nil {
		return "", 0, err
	}
	return token.String(), int64(options.AccessTokenExpireIn), nil
}

func VerifyJWT(token string) (*Token, error) {
	verifier, err := jwt.NewVerifierPS(jwt.PS512, publicPrivateKey.PublicKey)
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

	if !newClaims.IsSubject(os.Getenv("APP_NAME")) {
		return nil, fmt.Errorf(responses.InvalidToken)
	}

	// verify claims as you wish
	var isValid = newClaims.IsValidAt(time.Now())
	if !isValid {
		return nil, fmt.Errorf(responses.TokenExpired)
	}
	return (*Token)(&newClaims), nil
}

type Token jwt.RegisteredClaims

func (t *Token) HasRole(roles ...string) bool {
	tokenRoles := t.Audience

	for _, role := range roles {
		for _, tokenRole := range tokenRoles {
			if role == tokenRole {
				return true
			}
		}
	}
	return false
}
