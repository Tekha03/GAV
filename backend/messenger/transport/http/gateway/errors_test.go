package gateway

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	apperrors "shared/app_errors"
	"testing"
)

func TestWriteErrorMapsCategoryAndBody(t *testing.T) {
	recorder := httptest.NewRecorder()
	writeError(recorder, apperrors.New(
		apperrors.ChatNotFound,
		"chat not found",
		apperrors.WithDetail("chat_id", "42"),
	))

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusNotFound)
	}

	var response ErrorResponse
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.Error.Code != apperrors.CodeChatNotFound {
		t.Fatalf("code = %q, want %q", response.Error.Code, apperrors.CodeChatNotFound)
	}
	if response.Error.Category != apperrors.CategoryNotFound {
		t.Fatalf("category = %q, want %q", response.Error.Category, apperrors.CategoryNotFound)
	}
	if response.Error.Details["chat_id"] != "42" {
		t.Fatalf("details = %#v", response.Error.Details)
	}
}

func TestWriteErrorMasksInternalCause(t *testing.T) {
	recorder := httptest.NewRecorder()
	writeError(recorder, errors.New("database password leaked"))

	var response ErrorResponse
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.Error.Message != "internal server error" {
		t.Fatalf("message = %q", response.Error.Message)
	}
	if response.Error.Code != apperrors.CodeInternal {
		t.Fatalf("code = %q, want %q", response.Error.Code, apperrors.CodeInternal)
	}
}
