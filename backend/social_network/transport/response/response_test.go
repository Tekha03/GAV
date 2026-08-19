package response

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	apperrors "shared/app_errors"
	"testing"
)

func TestErrorMapsCategoryAndBody(t *testing.T) {
	recorder := httptest.NewRecorder()
	Error(recorder, apperrors.New(
		apperrors.PostNotFound,
		"post not found",
		apperrors.WithDetail("post_id", "42"),
	))

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusNotFound)
	}

	var body ErrorResponse
	if err := json.NewDecoder(recorder.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Error.Code != apperrors.CodePostNotFound {
		t.Fatalf("code = %q, want %q", body.Error.Code, apperrors.CodePostNotFound)
	}
	if body.Error.Category != apperrors.CategoryNotFound {
		t.Fatalf("category = %q, want %q", body.Error.Category, apperrors.CategoryNotFound)
	}
	if body.Error.Details["post_id"] != "42" {
		t.Fatalf("details = %#v", body.Error.Details)
	}
}

func TestErrorMasksUnknownError(t *testing.T) {
	recorder := httptest.NewRecorder()
	Error(recorder, errors.New("secret database failure"))

	var body ErrorResponse
	if err := json.NewDecoder(recorder.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusInternalServerError)
	}
	if body.Error.Code != apperrors.CodeInternal || body.Error.Message != "internal server error" {
		t.Fatalf("error = %#v", body.Error)
	}
}
