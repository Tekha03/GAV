package follow

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestNewFollow(t *testing.T) {
	followerID := uuid.New()
	followingID := uuid.New()

	t.Run("success", func(t *testing.T) {
		f, err := NewFollow(followerID, followingID)

		require.NoError(t, err)
		require.NotNil(t, f)
		require.Equal(t, followerID, f.FollowerID)
		require.Equal(t, followingID, f.FollowingID)
	})

	t.Run("follower is nil", func(t *testing.T) {
		f, err := NewFollow(uuid.Nil, followingID)

		require.ErrorIs(t, err, ErrFollowerIDNil)
		require.Nil(t, f)
	})

	t.Run("following is nil", func(t *testing.T) {
		f, err := NewFollow(followerID, uuid.Nil)

		require.ErrorIs(t, err, ErrFollowingIDNil)
		require.Nil(t, f)
	})
}
