package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const tokenIssuer = "chirpy-access"

func HashPassword(pw string) (string, error) {

	hashed_pw, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed_pw), nil

}

func ComparePasswords(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {

	authToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    tokenIssuer,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:   userID.String(),
	})

	return authToken.SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {

	customClaims := jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(tokenString, &customClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	if !token.Valid {
		return uuid.Nil, errors.New("token invalid")
	}
	userId := customClaims.Subject
	//issuer := customClaims.Issuer

	id, err := uuid.Parse(userId)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID: %w", err)
	}
	return id, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	rawBearerToken := headers.Get("Authorization")
	if rawBearerToken == "" {
		return "invalid bearer token", errors.New("no bearer token found")
	}

	bearerToken := strings.TrimPrefix(rawBearerToken, "Bearer ")
	if bearerToken == "" {
		return "", errors.New("bearer token does not exist")
	}

	return bearerToken, nil
}

func MakeRefreshToken() (string, error) {

	randNum, err := rand.Read(make([]byte, 32))
	if err != nil {
		return "", err
	}

	hexNum := hex.EncodeToString(make([]byte, randNum))

	return hexNum, nil
}
