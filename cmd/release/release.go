package main

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	"os"
)

func main() {
	dbUrl := os.Getenv("DATABASE_URL")
	if len(dbUrl) == 0 {
		log.Fatalf("DATABASE_URL undefined")
	}

	log.Printf("running migrations")
	m, err := migrate.New("file://migrate", dbUrl)
	if err != nil {
		log.Fatalf("failed to init migrate: %v", err)
	}

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Printf("no migrations performed (up to date)")
		} else {
			log.Fatalf("failed to up migrate: %v", err)
		}
	}
	log.Printf("ran migrations")
}
