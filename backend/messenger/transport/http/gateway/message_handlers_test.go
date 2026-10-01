package gateway

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"messenger/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type messageServiceStub struct {
	service.Service
	deletedBy  uuid.UUID
	deletedID  uuid.UUID
	readChatID uuid.UUID
	readUserID uuid.UUID
	err        error
}

func (s *messageServiceStub) DeleteMessage(_ context.Context, requesterID, messageID uuid.UUID) error {
	s.deletedBy, s.deletedID = requesterID, messageID
	return s.err
}

func (s *messageServiceStub) MarkAsRead(_ context.Context, chatID, requesterID uuid.UUID) error {
	s.readChatID, s.readUserID = chatID, requesterID
	return s.err
}

func requestWithParams(method, path string, userID *uuid.UUID, params map[string]string, body string) *http.Request {
	request := httptest.NewRequest(method, path, strings.NewReader(body))
	routeContext := chi.NewRouteContext()
	for key, value := range params {
		routeContext.URLParams.Add(key, value)
	}
	ctx := context.WithValue(request.Context(), chi.RouteCtxKey, routeContext)
	if userID != nil {
		ctx = context.WithValue(ctx, userIDContextKey, *userID)
	}
	return request.WithContext(ctx)
}

func TestDeleteMessageHandler(t *testing.T) {
	userID, messageID := uuid.New(), uuid.New()
	stub := &messageServiceStub{}
	recorder := httptest.NewRecorder()
	deleteMessage(stub).ServeHTTP(recorder, requestWithParams(
		http.MethodDelete,
		"/api/v1/messages/"+messageID.String(),
		&userID,
		map[string]string{"message_id": messageID.String()},
		"",
	))

	if recorder.Code != http.StatusNoContent || stub.deletedBy != userID || stub.deletedID != messageID {
		t.Fatalf("status=%d delete=(%s,%s)", recorder.Code, stub.deletedBy, stub.deletedID)
	}
}

func TestDeleteMessageHandlerRejectsInvalidOrUnauthenticatedRequests(t *testing.T) {
	stub := &messageServiceStub{}
	userID := uuid.New()
	for _, test := range []struct {
		name   string
		userID *uuid.UUID
		id     string
		status int
	}{
		{name: "invalid id", userID: &userID, id: "invalid", status: http.StatusBadRequest},
		{name: "missing auth", id: uuid.NewString(), status: http.StatusUnauthorized},
	} {
		t.Run(test.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			deleteMessage(stub).ServeHTTP(recorder, requestWithParams(
				http.MethodDelete,
				"/api/v1/messages/"+test.id,
				test.userID,
				map[string]string{"message_id": test.id},
				"",
			))
			if recorder.Code != test.status {
				t.Fatalf("status=%d, want %d", recorder.Code, test.status)
			}
		})
	}
}

func TestDeleteMessageHandlerMapsServiceError(t *testing.T) {
	userID, messageID := uuid.New(), uuid.New()
	stub := &messageServiceStub{err: errors.New("storage failed")}
	recorder := httptest.NewRecorder()
	deleteMessage(stub).ServeHTTP(recorder, requestWithParams(
		http.MethodDelete, "/", &userID,
		map[string]string{"message_id": messageID.String()}, "",
	))
	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("status=%d", recorder.Code)
	}
}

func TestMarkAsReadHandler(t *testing.T) {
	userID, chatID := uuid.New(), uuid.New()
	stub := &messageServiceStub{}
	recorder := httptest.NewRecorder()
	markAsRead(stub).ServeHTTP(recorder, requestWithParams(
		http.MethodPost, "/", &userID,
		map[string]string{"chat_id": chatID.String()},
		`{"user_id":"`+userID.String()+`"}`,
	))
	if recorder.Code != http.StatusOK || stub.readChatID != chatID || stub.readUserID != userID {
		t.Fatalf("status=%d read=(%s,%s)", recorder.Code, stub.readChatID, stub.readUserID)
	}
}

func TestMarkAsReadHandlerRejectsAnotherUser(t *testing.T) {
	userID, chatID := uuid.New(), uuid.New()
	recorder := httptest.NewRecorder()
	markAsRead(&messageServiceStub{}).ServeHTTP(recorder, requestWithParams(
		http.MethodPost, "/", &userID,
		map[string]string{"chat_id": chatID.String()},
		`{"user_id":"`+uuid.NewString()+`"}`,
	))
	if recorder.Code != http.StatusForbidden {
		t.Fatalf("status=%d, want %d", recorder.Code, http.StatusForbidden)
	}
}
