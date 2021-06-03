package utils

import (
	"errors"
	"strings"
)

func GetDomain(email string) (string, error) {
	emailParts := strings.Split(email, "@")
	if len(emailParts) < 2 {
		return "", errors.New("email is invalid")
	}
	return emailParts[1], nil
}
