package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func setupConfig() JWTConfig {
	return JWTConfig{
		Secret: []byte("secret"),
		TTL: time.Hour,
	}
}

func TestGenerateAndParseToken_Success(t *testing.T) {
	cfg := setupConfig()
	userID := uuid.New()

	token, err := GenerateAccessToken(userID, "user", cfg)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := ParseToken(token, cfg)

	assert.NoError(t, err)
	assert.NotNil(t, claims)

	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, "user", claims.Role)
}

func TestGenerateAccessToken_InvalidUserID(t *testing.T) {
	cfg := setupConfig()

	token, err := GenerateAccessToken(uuid.Nil, "user", cfg)

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Equal(t, ErrUserIDNil, err)
}

func TestGenerateAccessToken_EmptyRole(t *testing.T) {
	cfg := setupConfig()

	token, err := GenerateAccessToken(uuid.New(), "", cfg)

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Equal(t, ErrEmptyRole, err)
}

func TestGenerateAccessToken_ClaimsFields(t *testing.T) {
	cfg := setupConfig()
	userID := uuid.New()

	tokenStr, err := GenerateAccessToken(userID, "admin", cfg)
	assert.NoError(t, err)

	claims, err := ParseToken(tokenStr, cfg)
	assert.NoError(t, err)

	assert.Equal(t, "admin", claims.Role)
	assert.Equal(t, userID, claims.UserID)

	assert.NotNil(t, claims.ExpiresAt)
	assert.NotNil(t, claims.IssuedAt)

	assert.True(t, claims.ExpiresAt.After(claims.IssuedAt.Time))
}

func TestParseToken_WrongSecret(t *testing.T) {
	cfg := setupConfig()
	userID := uuid.New()

	token, err := GenerateAccessToken(userID, "user", cfg)
	assert.NoError(t, err)

	wrongConfig := JWTConfig{
		Secret: []byte("secret2"),
	}

	claims, err := ParseToken(token, wrongConfig)

	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.Equal(t, ErrInvalidToken, err)
}

func TestParseToken_InvalidTokenString(t *testing.T) {
	cfg := setupConfig()

	claims, err := ParseToken("invalid.token.string", cfg)

	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.Equal(t, ErrInvalidToken, err)
}

func TestParseToken_ExpiredToken(t *testing.T) {
	cfg := setupConfig()
	cfg.TTL = time.Millisecond
	userID := uuid.New()

	token, err := GenerateAccessToken(userID, "user", cfg)
	assert.NoError(t, err)

	time.Sleep(10 * time.Millisecond)
	claims, err := ParseToken(token, cfg)

	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.Equal(t, ErrInvalidToken, err)
}

func TestParseToken_UnexpectedSigningMethod(t *testing.T) {
	cfg := setupConfig()

	claims := &Claims{
		UserID: uuid.New(),
		Role: "user",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	tokenStr, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
	assert.NoError(t, err)

	parsedClaims, err := ParseToken(tokenStr, cfg)
	assert.Error(t, err)
	assert.Nil(t, parsedClaims)
}