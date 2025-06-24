package auth

import (
	"golang.org/x/crypto/bcrypt"
)

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
