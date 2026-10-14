package model

import (
	"context"
	"os"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// db is the PostgreSQL connection pool. Connect opens it once at startup.
var db *pgxpool.Pool

// Connect opens the database. The password is not an environment variable:
// it is read from the file DB_PASSWORD_FILE points at, which Docker Compose
// mounts from secrets/db_password.txt.
func Connect(ctx context.Context) error {
	cfg, err := pgxpool.ParseConfig("host=" + env("DB_HOST", "localhost") +
		" user=" + env("DB_USER", "alumni") +
		" dbname=" + env("DB_NAME", "alumni") +
		" sslmode=disable")
	if err != nil {
		return err
	}

	password, err := os.ReadFile(env("DB_PASSWORD_FILE", "secrets/db_password.txt"))
	if err != nil {
		return err
	}
	cfg.ConnConfig.Password = strings.TrimSpace(string(password))

	db, err = pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return err
	}
	return db.Ping(ctx)
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// Truncate empties every table and starts the ids from 1 again. Only the
// seed command uses it, for `seed -fresh`.
func Truncate(ctx context.Context) error {
	_, err := db.Exec(ctx, "TRUNCATE users, announcements, alumni RESTART IDENTITY")
	return err
}
