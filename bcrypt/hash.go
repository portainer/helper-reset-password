package bcrypt

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(plainTextPassword string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", nil
	}

	return string(hash), nil
}
