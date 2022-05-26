package main

import (
	"errors"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sethetter/go-web-starter/pkg/config"
)

func main() {
	c, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("failed to LoadConfig:", err)
	}

	m, err := migrate.New("file://sql", c.DatabaseURL)
	if err != nil {
		log.Fatalln("failed to migrate.New:", err)
	}

	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			log.Fatalln("failed to migrate Up:", err)
		}

		log.Println("no new migrations detected, schema is current")
	}
}
