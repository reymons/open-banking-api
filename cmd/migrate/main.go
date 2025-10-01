package main

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	"net/url"
	"os"
	"path/filepath"
)

func main() {
	var dbUrl = fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		url.QueryEscape(os.Getenv("DB_PASSWORD")),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	dir, err := filepath.Abs(os.Getenv("MIGRATIONS_DIR"))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Loading migrations from %s\n", dir)

	m, err := migrate.New("file://"+dir, dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		log.Fatal(err)
	}
	log.Println("Successfully migrated!")
}
