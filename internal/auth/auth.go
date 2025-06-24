package auth

import "golang.org/x/crypto/bcrypt"

func HashPassword(pw string) (string, error) {

	hashed_pw, err := bcrypt.GenerateFromPassword([]byte(pw), len(pw))
	if err != nil {
		return "", err
	}
	return string(hashed_pw), nil

}

func ComparePasswords(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
