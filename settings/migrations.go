package settings

import (
	"errors"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

func Migrations(db *sqlx.DB) error {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return errors.New("could not create postgres driver: " + err.Error())
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		return errors.New("could not create migrate instance: " + err.Error())
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return errors.New("could not run migrations:" + err.Error())
	}

	log.Println("Migrations ran successfully")
	return nil

}
