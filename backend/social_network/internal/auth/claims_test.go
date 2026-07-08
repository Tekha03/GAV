package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewClaims(t *testing.T) {
	tests := []struct {
		name      string
		userID    uuid.UUID
		role      string
		ttl       time.Duration
		wantError error
	}{
		{
			name:      "success",
			userID:    uuid.New(),
			role:      "user",
			ttl:       time.Hour,
			wantError: nil,
		},
		{
			name:      "empty user id",
			userID:    uuid.Nil,
			role:      "user",
			ttl:       time.Hour,
			wantError: ErrUserIDNil,
		},
		{
			name:      "empty role",
			userID:    uuid.New(),
			role:      "",
			ttl:       time.Hour,
			wantError: ErrEmptyRole,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			claims, err := NewClaims(test.userID, test.role, test.ttl)

			if test.wantError != nil {
				assert.Error(t, err)
				assert.Nil(t, claims)
				assert.Equal(t, test.wantError, err)

				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, claims)

			assert.Equal(t, test.userID, claims.UserID)
			assert.Equal(t, test.role, claims.Role)

			assert.NotNil(t, claims.ExpiresAt)
			assert.NotNil(t, claims.IssuedAt)

			assert.True(t, claims.ExpiresAt.After(claims.IssuedAt.Time))
		})
	}
}
