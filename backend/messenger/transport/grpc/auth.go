package grpc

import (
	"context"
	apperrors "shared/app_errors"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	grpcpkg "google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
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
			return nil, toGRPCError(err)
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
			return toGRPCError(err)
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
		return uuid.Nil, apperrors.New(apperrors.AuthTokenMissing, "authorization metadata is missing")
	}

	values := md.Get("authorization")
	if len(values) == 0 {
		return uuid.Nil, apperrors.New(apperrors.AuthTokenMissing, "authorization token is required")
	}

	return parseBearerToken(values[0], secret)
}

func parseBearerToken(header string, secret []byte) (uuid.UUID, error) {
	header = strings.TrimSpace(header)
	parts := strings.Fields(header)

	if len(parts) != 2 || parts[0] != "Bearer" || strings.TrimSpace(parts[1]) == "" {
		return uuid.Nil, apperrors.New(apperrors.AuthTokenInvalid, "invalid authorization token")
	}

	claims := &authClaims{}
	token, err := jwt.ParseWithClaims(parts[1], claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, apperrors.New(apperrors.AuthTokenInvalid, "invalid authorization token")
		}

		return secret, nil
	})

	if err != nil {
		return uuid.Nil, apperrors.Wrap(apperrors.AuthTokenInvalid, "invalid authorization token", err)
	}
	if token == nil || !token.Valid {
		return uuid.Nil, apperrors.New(apperrors.AuthTokenInvalid, "invalid authorization token")
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil || userID == uuid.Nil {
		return uuid.Nil, apperrors.New(apperrors.AuthTokenInvalid, "invalid authorization token")
	}

	return userID, nil
}
