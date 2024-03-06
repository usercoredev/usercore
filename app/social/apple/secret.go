package apple

// TODO: Will complete this later
func GenerateClientSecret(teamID, clientID, keyID string) (string, error) {
	/*
			privKey := cipher.ApplePrivateKey()

		// Create the Claims
		now := time.Now()
		claims := &jwt.StandardClaims{
			Issuer:    teamID,
			IssuedAt:  now.Unix(),
			ExpiresAt: now.Add(time.Hour*24*180 - time.Second).Unix(), // 180 days
			Audience:  "https://appleid.apple.com",
			Subject:   clientID,
		}

		token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
		token.Header["alg"] = "ES256"
		token.Header["kid"] = keyID

		return token.SignedString(privKey)
	*/
	return "", nil
}
