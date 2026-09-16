package testdb

import (
	"fmt"
	"net/url"
	"os"
	"testing"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DSN gives each integration test an isolated PostgreSQL schema.
func DSN(t *testing.T) string {
	t.Helper()
	base := os.Getenv("GAV_TEST_POSTGRES_DSN")
	if base == "" {
		t.Skip("GAV_TEST_POSTGRES_DSN is not set")
	}
	admin, err := gorm.Open(postgres.Open(base), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	schema := "test_" + uuid.NewString()[:8]
	if err := admin.Exec(fmt.Sprintf(`CREATE SCHEMA "%s"`, schema)).Error; err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := admin.Exec(fmt.Sprintf(`DROP SCHEMA "%s" CASCADE`, schema)).Error; err != nil {
			t.Error(err)
		}
		sqlDB, _ := admin.DB()
		_ = sqlDB.Close()
	})
	u, err := url.Parse(base)
	if err != nil {
		t.Fatal(err)
	}
	query := u.Query()
	query.Set("search_path", schema)
	u.RawQuery = query.Encode()
	return u.String()
}

func Open(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(postgres.Open(DSN(t)), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		sqlDB, _ := db.DB()
		_ = sqlDB.Close()
	})
	return db
}
