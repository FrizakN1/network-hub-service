package database

import (
	"backend/settings"
	"backend/utils"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

var Link *sql.DB

func Connection(config *settings.Setting) {
	var e error
	Link, e = sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DbHost,
		config.DbPort,
		config.DbUser,
		config.DbPass,
		config.DbName))
	if e != nil {
		fmt.Println(e)
		utils.Logger.Println(e)
		return
	}

	e = Link.Ping()
	if e != nil {
		fmt.Println(e)
		utils.Logger.Println(e)
		return
	}

	if e = goose.SetDialect("postgres"); e != nil {
		fmt.Println(e)
		utils.Logger.Println(e)
		return
	}

	if e = goose.Up(Link, "migrations"); e != nil {
		fmt.Println(e)
		utils.Logger.Println(e)
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

	if len(errorsList) > 0 {
		for _, i := range errorsList {
			fmt.Println(i)
			utils.Logger.Println(i)
		}
	}

	LoadAddressElementTypeMap(addressElementTypeMap)
	LoadSession(sessionMap)
}
