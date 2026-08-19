package grpc

import (
	"errors"
	apperrors "shared/app_errors"
	"testing"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestToGRPCErrorMapsCategory(t *testing.T) {
	err := toGRPCError(apperrors.New(apperrors.ChatAccessDenied, "chat access denied"))
	if got := status.Code(err); got != codes.PermissionDenied {
		t.Fatalf("code = %s, want %s", got, codes.PermissionDenied)
	}
	if got := status.Convert(err).Message(); got != "chat access denied" {
		t.Fatalf("message = %q", got)
	}
	details := status.Convert(err).Details()
	if len(details) != 1 {
		t.Fatalf("details count = %d, want 1", len(details))
	}
	info, ok := details[0].(*errdetails.ErrorInfo)
	if !ok || info.Reason != string(apperrors.CodeChatAccessDenied) {
		t.Fatalf("error info = %#v", details[0])
	}
}

func TestToGRPCErrorMasksInternalCause(t *testing.T) {
	err := toGRPCError(errors.New("database password leaked"))
	if got := status.Code(err); got != codes.Internal {
		t.Fatalf("code = %s, want %s", got, codes.Internal)
	}
	if got := status.Convert(err).Message(); got != "internal server error" {
		t.Fatalf("message = %q", got)
	}
}
