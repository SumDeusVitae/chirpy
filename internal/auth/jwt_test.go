package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMakeJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "my_secret_key"
	expiresIn := time.Hour

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Validate the token
	parsedUserID, err := ValidateJWT(token, tokenSecret)
	assert.NoError(t, err)
	assert.Equal(t, userID, parsedUserID)
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "my_secret_key"
	expiresIn := time.Hour

	// Create a token
	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	assert.NoError(t, err)

	// Validate the token
	parsedUserID, err := ValidateJWT(token, tokenSecret)
	assert.NoError(t, err)
	assert.Equal(t, userID, parsedUserID)

	// Test with an invalid token
	invalidToken := "invalid.token.string"
	_, err = ValidateJWT(invalidToken, tokenSecret)
	assert.Error(t, err)

	// Test with a token that has an invalid signing method
	_, err = ValidateJWT(token, "wrong_secret_key")
	assert.Error(t, err)
}
