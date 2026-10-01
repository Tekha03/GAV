package postgres

import (
	"context"
	"fmt"
	"os"
	"testing"

	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestMain(m *testing.M) {
	os.Exit(runTestsWithPostgres(m))
}

func runTestsWithPostgres(m *testing.M) int {
	if os.Getenv("GAV_TEST_POSTGRES_DSN") != "" {
		return m.Run()
	}

	ctx := context.Background()
	container, err := tcpostgres.Run(ctx, "postgres:15-alpine",
		tcpostgres.WithDatabase("gav_test"),
		tcpostgres.WithUsername("gav_test"),
		tcpostgres.WithPassword("gav_test"),
		tcpostgres.BasicWaitStrategies(),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "testcontainers: PostgreSQL unavailable, database tests will be skipped: %v\n", err)
		return m.Run()
	}

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		_ = container.Terminate(ctx)
		fmt.Fprintf(os.Stderr, "testcontainers: get PostgreSQL DSN: %v\n", err)
		return 1
	}
	if err := os.Setenv("GAV_TEST_POSTGRES_DSN", dsn); err != nil {
		_ = container.Terminate(ctx)
		fmt.Fprintf(os.Stderr, "testcontainers: set PostgreSQL DSN: %v\n", err)
		return 1
	}

	code := m.Run()
	if err := container.Terminate(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "testcontainers: terminate PostgreSQL: %v\n", err)
		if code == 0 {
			code = 1
		}
	}
	return code
}
