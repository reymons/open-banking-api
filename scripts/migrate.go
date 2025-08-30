package main 

import (
	"log"
    "os"
    "path/filepath"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
    dir, err := filepath.Abs("db", "migrations")
    if err != nil {
        log.Fatal(err)
    }
	m, err := migrate.New("file://" + dir, os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		log.Fatal(err)
	}
    log.Println("Successfully migrated!")
}
