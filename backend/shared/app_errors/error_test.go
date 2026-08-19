package apperrors

import (
	"context"
	"errors"
	"fmt"
	"testing"
)

func TestWrapPreservesCauseAndMetadata(t *testing.T) {
	cause := errors.New("database failure")
	err := Wrap(UserNotFound, "user not found", cause, WithDetail("user_id", "42"))

	if !errors.Is(err, cause) {
		t.Fatal("wrapped error does not preserve its cause")
	}

	code, ok := CodeOf(fmt.Errorf("service: %w", err))
	if !ok || code != CodeUserNotFound {
		t.Fatalf("CodeOf() = %q, %v; want %q, true", code, ok, CodeUserNotFound)
	}

	category, ok := CategoryOf(err)
	if !ok || category != CategoryNotFound {
		t.Fatalf("CategoryOf() = %q, %v; want %q, true", category, ok, CategoryNotFound)
	}

	details, ok := DetailsOf(err)
	if !ok || details["user_id"] != "42" {
		t.Fatalf("DetailsOf() = %#v, %v", details, ok)
	}

	details["user_id"] = "changed"
	original, _ := DetailsOf(err)
	if original["user_id"] != "42" {
		t.Fatal("DetailsOf returned the original mutable map")
	}
}

func TestErrorsIsComparesCodes(t *testing.T) {
	err := New(UserNotFound, "first message")
	target := New(UserNotFound, "another message")

	if !errors.Is(err, target) {
		t.Fatal("errors with the same code must match")
	}
}

func TestWrapNilReturnsNil(t *testing.T) {
	if err := Wrap(Internal, "failed", nil); err != nil {
		t.Fatalf("Wrap() = %v; want nil", err)
	}
}

func TestNormalize(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantCode Code
	}{
		{name: "cancelled", err: context.Canceled, wantCode: CodeRequestCancelled},
		{name: "timeout", err: context.DeadlineExceeded, wantCode: CodeRequestTimeout},
		{name: "unknown", err: errors.New("secret database error"), wantCode: CodeInternal},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			normalized := Normalize(test.err)
			if normalized.Code != test.wantCode {
				t.Fatalf("Normalize().Code = %q; want %q", normalized.Code, test.wantCode)
			}
			if !errors.Is(normalized, test.err) {
				t.Fatal("Normalize did not preserve the original cause")
			}
		})
	}
}
