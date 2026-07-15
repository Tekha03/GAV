package grpc

import (
	"context"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	grpcpkg "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type authContextKey struct{}

var userIDContextKey = authContextKey{}

type authClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

type authServerStream struct {
	grpcpkg.ServerStream
	ctx context.Context
}

func AuthUnaryInterceptor(jwtSecret string) grpcpkg.UnaryServerInterceptor {
	secret := []byte(jwtSecret)

	return func(
		ctx context.Context,
		req any,
		info *grpcpkg.UnaryServerInfo,
		handler grpcpkg.UnaryHandler,
	) (any, error) {
		userID, err := userIDFromMetadata(ctx, secret)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "unauthenticated")
		}

		ctx = context.WithValue(ctx, userIDContextKey, userID)
		return handler(ctx, req)
	}
}

func AuthStreamInterceptor(jwtSecret string) grpcpkg.StreamServerInterceptor {
	secret := []byte(jwtSecret)

	return func(
		srv any,
		stream grpcpkg.ServerStream,
		info *grpcpkg.StreamServerInfo,
		handler grpcpkg.StreamHandler,
	) error {
		userID, err := userIDFromMetadata(stream.Context(), secret)
		if err != nil {
			return status.Error(codes.Unauthenticated, "unauthenticated")
		}

		ctx := context.WithValue(stream.Context(), userIDContextKey, userID)
		return handler(srv, &authServerStream{
			ServerStream: stream,
			ctx:          ctx,
		})
	}
}

func (s *authServerStream) Context() context.Context {
	return s.ctx
}

func CurrentUserID(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(userIDContextKey).(uuid.UUID)
	return id, ok
}

func userIDFromMetadata(ctx context.Context, secret []byte) (uuid.UUID, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return uuid.Nil, ErrMissingMetadata
	}

	values := md.Get("authorization")
	if len(values) == 0 {
		return uuid.Nil, ErrMissingAuthMetadata
	}

	return parseBearerToken(values[0], secret)
}

func parseBearerToken(header string, secret []byte) (uuid.UUID, error) {
	header = strings.TrimSpace(header)
	parts := strings.Fields(header)

	if len(parts) != 2 || parts[0] != "Bearer" || strings.TrimSpace(parts[1]) == "" {
		return uuid.Nil, ErrIvalidAuthorization
	}

	claims := &authClaims{}
	token, err := jwt.ParseWithClaims(parts[1], claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSignMethod
		}

		return secret, nil
	})

	if err != nil || !token.Valid {
		return uuid.Nil, ErrInvalidToken
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil || userID == uuid.Nil {
		return uuid.Nil, ErrInvalidSubject
	}

	return userID, nil
}
