package database

import (
	"backend/utils"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"os"
)

var Link *sql.DB

func Connection() {
	var err error
	Link, err = sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME")))
	if err != nil {
		fmt.Println(err)
		utils.Logger.Println(err)
		return
	}

	if err = Link.Ping(); err != nil {
		utils.Logger.Println(err)
		return
	}

	if err = goose.SetDialect("postgres"); err != nil {
		utils.Logger.Println(err)
		return
	}

	if err = goose.Up(Link, "migrations"); err != nil {
		utils.Logger.Println(err)
		return
	}

	errorsList := make([]string, 0)

	errorsList = append(errorsList, prepareReferences()...)
	errorsList = append(errorsList, prepareHouse()...)
	errorsList = append(errorsList, prepareUsers()...)
	errorsList = append(errorsList, prepareNodes()...)
	errorsList = append(errorsList, prepareHardware()...)
	errorsList = append(errorsList, prepareFile()...)
	errorsList = append(errorsList, prepareSwitch()...)
	errorsList = append(errorsList, prepareEvent()...)

	if len(errorsList) > 0 {
		for _, i := range errorsList {
			fmt.Println(i)
			utils.Logger.Println(i)
		}
	}

	LoadAddressElementTypeMap(addressElementTypeMap)
	LoadSession(sessionMap)
}
