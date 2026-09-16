package grpc

import (
	"context"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestAuthUnaryInterceptorAcceptsSocialNetworkJWT(t *testing.T) {
	const secret = "shared-secret"
	userID := uuid.New()
	token := signedToken(t, secret, userID, time.Now().Add(time.Hour))
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "Bearer "+token))

	var captured uuid.UUID
	_, err := AuthUnaryInterceptor(secret)(ctx, nil, nil, func(ctx context.Context, _ any) (any, error) {
		captured, _ = CurrentUserID(ctx)
		return struct{}{}, nil
	})
	if err != nil {
		t.Fatalf("interceptor returned error: %v", err)
	}
	if captured != userID {
		t.Fatalf("authenticated user = %s, want %s", captured, userID)
	}
}

func TestAuthUnaryInterceptorRejectsExpiredAndInvalidJWT(t *testing.T) {
	const secret = "shared-secret"
	tests := map[string]string{
		"expired":      signedToken(t, secret, uuid.New(), time.Now().Add(-time.Minute)),
		"wrong secret": signedToken(t, "different-secret", uuid.New(), time.Now().Add(time.Hour)),
		"malformed":    "not-a-jwt",
	}

	for name, token := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "Bearer "+token))
			_, err := AuthUnaryInterceptor(secret)(ctx, nil, nil, func(context.Context, any) (any, error) {
				t.Fatal("handler must not be called")
				return nil, nil
			})
			if status.Code(err) != codes.Unauthenticated {
				t.Fatalf("code = %s, want %s", status.Code(err), codes.Unauthenticated)
			}
		})
	}
}

func signedToken(t *testing.T, secret string, userID uuid.UUID, expiresAt time.Time) string {
	t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  userID.String(),
		"role": "user",
		"exp":  expiresAt.Unix(),
		"iat":  time.Now().Unix(),
	})
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatal(err)
	}
	return signed
}
