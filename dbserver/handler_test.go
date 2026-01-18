package dbserver

import (
	"bytes"
	"encoding/json"
	"gav/storage/memory"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupTestServer() *gin.Engine {
	r := gin.Default()

	userRepo := memory.NewUserRepository()
	authRepo := NewAuthHandler(userRepo)

	r.POST("/register", authRepo.Register)
	r.POST("/login", authRepo.Login)

	return r
}

type Credentials struct {
	Email 	 string	`json:"email"`
	Password string	`json:"password"`
}

func TestRegister_Success(t *testing.T) {
	r := setupTestServer()

	body := Credentials{
		Email: 		"test@example.com",
		Password: 	"superPassword",
	}
	data, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/register", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestRegister_DuplicateEmail(t *testing.T) {
	r := setupTestServer()

	body := Credentials{
		Email: 		"dup@example.com",
		Password: 	"pass123",
	}
	data, _ := json.Marshal(body)

	// Первый раз – ок
	req1, _ := http.NewRequest("POST", "/register", bytes.NewReader(data))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)

	// Второй раз – ошибка
	req2, _ := http.NewRequest("POST", "/register", bytes.NewReader(data))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	if w2.Code == http.StatusOK {
		t.Fatalf("expected error on duplicate email, got 200 OK")
	}
}
