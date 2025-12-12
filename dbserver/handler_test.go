package dbserver

import (
	"bytes"
	"encoding/json"
	"gav/user"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) {
	var err error
	DB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open in-memory DB: %v", err)
	}

	if err := DB.AutoMigrate(&user.User{}); err != nil {
		t.Fatalf("migration failed: %v", err)
	}
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/register", RegisterHandler)
	r.GET("users/:id", GetUserHandler)
	return r
}

func TestRegisterHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB(t)
	router := setupRouter()

	input := user.User{
		ID: 10,
	}
	body, _ := json.Marshal(input)

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var created user.User
	if err := json.Unmarshal(w.Body.Bytes(), &created); err != nil {
		t.Fatalf("response JSON unmarshal failed: %v", err)
	}

	if created.ID != 10 {
		t.Fatalf("expected user ID 10, got %d", created.ID)
	}

	var stored user.User
	if err := DB.First(&stored, 10).Error; err != nil {
		t.Fatalf("user was not inserted into DB: %v", err)
	}

}

func TestGetUserHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB(t)

	DB.Create(&user.User{ID: 25})
	router := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/users/25", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var returned user.User
	if err := json.Unmarshal(w.Body.Bytes(), &returned); err != nil {
		t.Fatalf("bad JSON: %v", err)
	}

	if returned.ID != 25 {
		t.Fatalf("expectd ID 25, got %d", returned.ID)
	}
}

func TestGetUserNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB(t)
	router := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/users/999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}
