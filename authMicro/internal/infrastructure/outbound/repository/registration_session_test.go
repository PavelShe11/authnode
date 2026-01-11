//go:build integration
// +build integration

package repository

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/PavelShe11/studbridge/authMicro/internal/entity"
	"github.com/PavelShe11/studbridge/authMicro/internal/port"
	trmsql "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	testDB        *sqlx.DB
	testContainer testcontainers.Container
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	// Setup PostgreSQL container
	req := testcontainers.ContainerRequest{
		Image:        "postgres:16-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
			"POSTGRES_DB":       "test_db",
		},
		WaitStrategy: wait.ForLog("database system is ready to accept connections").
			WithOccurrence(2).
			WithStartupTimeout(60 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to start container: %s", err))
	}

	testContainer = container

	// Get connection string
	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "5432")
	connStr := fmt.Sprintf("postgres://test:test@%s:%s/test_db?sslmode=disable", host, port.Port())

	// Connect to database
	testDB, err = sqlx.Connect("postgres", connStr)
	if err != nil {
		container.Terminate(ctx)
		panic(fmt.Sprintf("failed to connect to database: %s", err))
	}

	// Run migrations
	runMigrations(testDB)

	// Run tests
	code := m.Run()

	// Cleanup
	testDB.Close()
	container.Terminate(ctx)

	os.Exit(code)
}

func runMigrations(db *sqlx.DB) {
	migration := `
	CREATE TABLE IF NOT EXISTS registration_session (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		code TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		code_expires TIMESTAMP NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW()
	);
	CREATE INDEX IF NOT EXISTS idx_registration_session_email ON registration_session(email);
	CREATE INDEX IF NOT EXISTS idx_registration_session_code_expires ON registration_session(code_expires);
	`
	_, err := db.Exec(migration)
	if err != nil {
		panic(fmt.Sprintf("failed to run migrations: %s", err))
	}
}

// setupTest creates a fresh repository with clean database
func setupTest(t *testing.T) port.RegistrationSessionRepository {
	_, err := testDB.Exec("TRUNCATE TABLE registration_session")
	require.NoError(t, err)

	getter := trmsql.NewCtxGetter(testDB)
	return NewRegistrationSessionRepository(testDB, getter)
}