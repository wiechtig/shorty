package testutil

import (
	"context"
	"fmt"
	"go.wiechtig.com/shorty/internal/shared"
	"go.wiechtig.com/shorty/internal/store"
	"go.wiechtig.com/shorty/pkg/nanoid"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// needs a running postgres instance
var databaseURLPrefix = "postgres://postgres:postgres@localhost:5432"

func WithDatabase[TB testing.TB](ctx context.Context, tb TB, test func(t TB, db *pgxpool.Pool, store *store.Queries)) {
	dbname, _ := nanoid.New(12)
	databaseURL := databaseURLPrefix + "/" + dbname
	tb.Logf("database name: %s", dbname)

	pool, err := pgxpool.New(ctx, databaseURLPrefix)
	if err != nil {
		tb.Fatalf("Unable to connect to database: %v\n", err)
		return
	}
	defer func() {
		pool.Close()
	}()
	err = createDatabase(ctx, pool, dbname)
	if err != nil {
		tb.Fatalf("Unable to create database: %v\n", err)
		return
	}
	defer func() {
		err := dropDatabase(ctx, pool, dbname)
		if err != nil {
			tb.Logf("Unable to drop database: %v", err)
		}
	}()

	db := shared.SetupDatabase(databaseURL)
	defer func() {
		db.Close()
	}()

	root, err := findProjectRoot()
	if err != nil {
		tb.Fatalf("Could not find project root: %v", err)
	}
	migrationPath := filepath.Join(root, "db", "migrations")
	shared.RunMigrations(db, migrationPath)

	dataPath := filepath.Join(root, "db", "data")
	err = InsertTestData(ctx, db, dataPath)
	if err != nil {
		tb.Fatalf("Unable to insert test data: %v\n", err)
		return
	}

	s := store.New(db)

	test(tb, db, s)
}

// createDatabase creates a new database with the specified name.
func createDatabase(ctx context.Context, db *pgxpool.Pool, name string) error {
	_, err := db.Exec(ctx, `CREATE DATABASE `+sanitizeDatabaseName(name)+`;`)
	return err
}

// dropDatabase drops the specific database.
func dropDatabase(ctx context.Context, db *pgxpool.Pool, name string) error {
	_, err := db.Exec(ctx, `DROP DATABASE `+sanitizeDatabaseName(name)+`;`)
	return err
}

// sanitizeDatabaseName ensures that the database name is a valid postgres identifier.
func sanitizeDatabaseName(schema string) string {
	return pgx.Identifier{schema}.Sanitize()
}

func InsertTestData(ctx context.Context, db *pgxpool.Pool, path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("failed to access path %s: %w", path, err)
	}

	if info.IsDir() {
		files, err := os.ReadDir(path)
		if err != nil {
			return fmt.Errorf("failed to read directory: %w", err)
		}

		for _, file := range files {
			if filepath.Ext(file.Name()) == ".sql" {
				err := executeSQLFile(ctx, db, filepath.Join(path, file.Name()))
				if err != nil {
					return err
				}
			}
		}
	} else if filepath.Ext(path) == ".sql" {
		err := executeSQLFile(ctx, db, path)
		if err != nil {
			return err
		}
	}
	return nil
}

func executeSQLFile(ctx context.Context, db *pgxpool.Pool, filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	queries := strings.Split(string(content), ";")
	for _, query := range queries {
		query = strings.TrimSpace(query)
		if query != "" {
			_, err := db.Exec(ctx, query)
			if err != nil {
				return fmt.Errorf("failed to execute query from file %s: %w", filePath, err)
			}
		}
	}

	return nil
}

func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", fmt.Errorf("project root not found")
}
