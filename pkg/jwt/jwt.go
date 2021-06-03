package jwt

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type JWT struct {
	secret    []byte
	issuer    string
	expiresAt time.Duration
}

type BasicClaims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func NewJWT(secret, issuer string, expiresAt time.Duration) JWT {
	return JWT{
		secret:    []byte(secret),
		issuer:    issuer,
		expiresAt: expiresAt,
	}
}

func (j *JWT) CreateJWT(email string) (string, error) {
	// Create the claims
	claims := BasicClaims{
		email,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(j.expiresAt).Unix(),
			Issuer:    j.issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(j.secret)
	if err != nil {
		return "", fmt.Errorf("fail to sign token: %v", err)
	}
	return signedToken, nil
}

func (j *JWT) ParseToken(inputToken string) (BasicClaims, error) {
	token, err := jwt.ParseWithClaims(inputToken,
		&BasicClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return j.secret, nil
		},
	)
	if err != nil {
		return BasicClaims{}, fmt.Errorf("fail to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*BasicClaims); ok && token.Valid {
		return *claims, nil
	}

	return BasicClaims{}, fmt.Errorf("fail to parse token: Unknown reason")
}
