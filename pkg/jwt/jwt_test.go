package jwt

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestJWT_CreateJWT(t *testing.T) {
	j := &JWT{
		secret:    []byte("secret"),
		issuer:    "issuer",
		expiresAt: 10 * time.Hour,
	}
	email := "test@coke.com.us"
	token, err := j.CreateJWT(email)
	assert.NoError(t, err)

	basicClaims, err := j.ParseToken(token)
	assert.NoError(t, err)
	assert.NotNil(t, basicClaims)
	assert.Equal(t, basicClaims.Email, email)
}
