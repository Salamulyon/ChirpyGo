package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {

	authToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now() + expiresIn,
		Subject:   string(userID)})

	signedAuthToken, err := authToken.SignedString(jwt.SigningMethodHS256)
	if err != nil {
		return "", err
	}

	return signedAuthToken, nil
}

func ValidateJWT() {

}
