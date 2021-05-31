package criptography

import "golang.org/x/crypto/argon2"

func HashPassword(password, salt string) []byte {
	hPassword := argon2.Key([]byte(password), []byte(salt), 3, 128, 1, 32)
	return hPassword
}
