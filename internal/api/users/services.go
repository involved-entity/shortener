package users

import (
	"crypto/rand"
	"errors"
	"log"
	"math/big"
)

func GenerateSecureToken() (string, error) {
	const digits = "0123456789"
	const length = 5
	token := make([]byte, length)

	for i := range token {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			log.Printf("failed to generate token: %v", err)
			return "", errors.New("failed to generate token")
		}
		token[i] = digits[num.Int64()]
	}

	return string(token), nil
}
