package main

import (
	"errors"
	"log"
	"os"
	"tracker-bot/internal/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
	if os.Getenv("HOST_DB") == "" {
		_ = godotenv.Load(".env")
	}

	cfg, err := config.ParseConfig()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	dsn := cfg.PostgresURL()

	m, err := migrate.New(
		"file://migrations",
		dsn,
	)
	if err != nil {
		log.Fatalf("migrate new: %v", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("no migrations to apply")
			return
		}
		log.Fatalf("migrate up: %v", err)
	}

	log.Println("migrations applied")
}
