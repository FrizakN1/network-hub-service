package database

import (
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"log"
)

func InitDatabase(migration string) (Database, error) {
	d := new(DefaultDatabase)

	if err := d.Connect(); err != nil {
		return nil, err
	}

	if err := goose.SetDialect("postgres"); err != nil {
		return nil, err
	}

	switch migration {
	case "up":
		if err := goose.Up(d.db, "migrations"); err != nil {
			return nil, err
		}
	case "down":
		if err := goose.Down(d.db, "migrations"); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("invalid migration direction: %s. Use 'up' or 'down'", migration)
	}
	if err := goose.Up(d.db, "migrations"); err != nil {
		return nil, err
	}

	errorsList := d.PrepareQuery()
	if len(errorsList) > 0 {
		for _, err := range errorsList {
			log.Println(err)
		}

		return nil, errors.New("failed to prepare query")
	}

	return d, nil
}
