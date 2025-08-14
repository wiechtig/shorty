package shared

import (
	"context"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"log/slog"
)

func SetupDatabase(databaseUrl string) *pgxpool.Pool {
	dbpool, err := pgxpool.New(context.Background(), databaseUrl)
	if err != nil {
		slog.Error("Unable to create connection pool", slog.Any("error", err))
		panic(err)
	}
	return dbpool
}

func RunMigrations(pool *pgxpool.Pool, dir string) {
	db := stdlib.OpenDBFromPool(pool)
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		slog.Error("Unable to connect to database", slog.Any("error", err))
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+dir,
		"postgres", driver)
	defer m.Close()
	if err != nil {
		slog.Error("Unable to create migrate instance", slog.Any("error", err))
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		slog.Error("Migration failed", slog.Any("error", err))
	} else {
		slog.Info("Database migrations completed successfully", slog.String("directory", dir))
	}
}
