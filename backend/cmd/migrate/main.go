// Command migrate applies every .sql file in ./migrations in filename order.
// Each file is expected to be idempotent (CREATE TABLE IF NOT EXISTS ...),
// which keeps the runner tiny — no version table needed for this MVP.
package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"sort"

	"travelcrm/internal/config"
	"travelcrm/internal/db"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	pool, err := db.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("connect: %v", err)
	}
	defer pool.Close()

	dir := migrationsDir()
	entries, err := filepath.Glob(filepath.Join(dir, "*.sql"))
	if err != nil {
		log.Fatalf("glob migrations: %v", err)
	}
	sort.Strings(entries)

	if len(entries) == 0 {
		log.Fatalf("no migration files found in %s", dir)
	}

	for _, file := range entries {
		sqlBytes, err := os.ReadFile(file)
		if err != nil {
			log.Fatalf("read %s: %v", file, err)
		}
		if _, err := pool.Exec(ctx, string(sqlBytes)); err != nil {
			log.Fatalf("apply %s: %v", filepath.Base(file), err)
		}
		log.Printf("applied %s", filepath.Base(file))
	}
	log.Println("migrations complete")
}

// migrationsDir resolves the migrations folder relative to the module root,
// so the command works whether run from ./backend or via `go run ./cmd/migrate`.
func migrationsDir() string {
	if dir := os.Getenv("MIGRATIONS_DIR"); dir != "" {
		return dir
	}
	return "migrations"
}
