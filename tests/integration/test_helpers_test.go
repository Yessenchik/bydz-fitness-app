package integration

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func openTestDB(t *testing.T) *sql.DB {
	t.Helper()

	dsn := os.Getenv("TEST_DATABASE_DSN")
	if dsn == "" {
		dsn = "host=localhost port=5432 user=gym password=gym_password dbname=membership_db sslmode=disable"
	}

	db, err := sql.Open("postgres", dsn)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	require.NoError(t, err)

	t.Cleanup(func() {
		_, _ = db.ExecContext(
			context.Background(),
			`
			DELETE FROM attendances;
			DELETE FROM subscriptions;
			DELETE FROM membership_plans
			WHERE name LIKE 'Test%';
			`,
		)

		_ = db.Close()
	})

	return db
}
