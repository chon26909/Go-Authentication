package security

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func EncryptPassword(password string) (string, error) {

	salt := bcrypt.DefaultCost

	fmt.Println("salt ", salt)

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), salt)

	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func VerifyPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
