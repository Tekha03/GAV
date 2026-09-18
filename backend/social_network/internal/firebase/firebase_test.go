package firebase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClientDoesNotSilentlyAcceptMissingCredentials(t *testing.T) {
	_, err := NewClient(context.Background(), "")
	require.Error(t, err)
	_, err = NewClient(context.Background(), t.TempDir()+"/missing.json")
	require.Error(t, err)
	require.ErrorIs(t, (&Client{}).SendPush(context.Background(), "token", "title", "body", nil), ErrDisabled)
}
