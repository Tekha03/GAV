package auth

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword_Success(t *testing.T) {
	hasher := &PasswordHasher{}
	password := "my-secure-password"

	hash, err := hasher.HashPassword(password)
	assert.NoError(t, err)

	assert.NotEmpty(t, hash)
	assert.NotEqual(t, password, hash)
}

func TestHashPassword_GenerateDifferentHashes(t *testing.T) {
	hasher := &PasswordHasher{}
	password := "same-password"

	hash1, err := hasher.HashPassword(password)
	assert.NoError(t, err)

	hash2, err := hasher.HashPassword(password)
	assert.NoError(t, err)

	assert.NotEqual(t, hash1, hash2)

	assert.True(t, hasher.CheckPassword(hash1, password))
	assert.True(t, hasher.CheckPassword(hash2, password))
}

func TestHashPassword_EmptyPassword(t *testing.T) {
	hasher := &PasswordHasher{}

	hash, err := hasher.HashPassword("")

	assert.NoError(t, err)
	assert.NotEmpty(t, hash)

	ok := hasher.CheckPassword(hash, "")
	assert.True(t, ok)
}

func TestHashPassword_LongPassword(t *testing.T) {
	hasher := &PasswordHasher{}

	longPassword := strings.Repeat("a", 100)
	hash, err := hasher.HashPassword(longPassword)

	assert.Error(t, err)
	assert.Empty(t, hash)
	assert.Equal(t, bcrypt.ErrPasswordTooLong, err)
}

func TestHashPassword_DifferentHash(t *testing.T) {
	hasher := &PasswordHasher{}

	hash1, err := hasher.HashPassword("password1")
	assert.NoError(t, err)

	ok := hasher.CheckPassword(hash1, "password2")
	assert.False(t, ok)
}

func TestCheckPassword_CorrectPassword(t *testing.T) {
	hasher := &PasswordHasher{}
	password := "password123"

	hash, err := hasher.HashPassword(password)
	assert.NoError(t, err)

	ok := hasher.CheckPassword(hash, password)
	assert.True(t, ok)
}

func TestCheckPassword_WrongPassword(t *testing.T) {
	hasher := &PasswordHasher{}

	hash, err := hasher.HashPassword("correct-password")
	assert.NoError(t, err)

	ok := hasher.CheckPassword(hash, "wrong-password")
	assert.False(t, ok)
}

func TestCheckPassword_InvalidHash(t *testing.T) {
	hasher := &PasswordHasher{}

	ok := hasher.CheckPassword("not-a-valid-hash", "password")
	assert.False(t, ok)
}
