package database

import (
	"backend/settings"
	"backend/utils"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
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

	errorsList := make([]string, 0)

	errorsList = append(errorsList, prepareEnums()...)
	errorsList = append(errorsList, prepareRequests()...)

	if len(errorsList) > 0 {
		for _, i := range errorsList {
			fmt.Println(i)
			utils.Logger.Println(i)
		}
	}

	LoadEnums(enumsMap)
}
