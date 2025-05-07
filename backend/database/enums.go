package database

import (
	"backend/utils"
	"database/sql"
	"fmt"
)

type Enum struct {
	ID             int
	Name           string
	Value          string
	TranslateValue string
	CreatedAt      int64
}

var query map[string]*sql.Stmt
var enumsMap map[string]map[string]bool

func prepareEnums() []string {
	var e error
	errorsList := make([]string, 0)

	enumsMap = make(map[string]map[string]bool)

	if query == nil {
		query = make(map[string]*sql.Stmt)
	}

	query["STREET_TYPE"], e = Link.Prepare(`SELECT * FROM "Street_type" ORDER BY id`)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["HOUSE_TYPE"], e = Link.Prepare(`SELECT * FROM "House_type" ORDER BY id`)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	return errorsList
}

func LoadEnums(m map[string]map[string]bool) {
	keys := [2]string{
		"STREET_TYPE",
		"HOUSE_TYPE",
	}

	for _, key := range keys {
		stmt, ok := query[key]
		if !ok {
			fmt.Println("ошибка загрузки перечислений, необходимо остановить работу сервера и обратиться к разработчику")
			utils.Logger.Println("ошибка загрузки перечислений, необходимо остановить работу сервера и обратиться к разработчику")
			return
		}

		rows, e := stmt.Query()
		if e != nil {
			fmt.Println(e)
			utils.Logger.Println(e)
			fmt.Println("ошибка загрузки перечислений, необходимо остановить работу сервера и обратиться к разработчику")
			utils.Logger.Println("ошибка загрузки перечислений, необходимо остановить работу сервера и обратиться к разработчику")
			return
		}

		if m[key] == nil {
			m[key] = make(map[string]bool)
		}

		for rows.Next() {
			var enum AddressElementType
			e = rows.Scan(
				&enum.ID,
				&enum.Name,
				&enum.ShortName,
			)
			if e != nil {
				utils.Logger.Println(e)
				fmt.Println(e)
				return
			}

			m[key][enum.Name] = true
			m[key][enum.ShortName] = true
		}

		_ = rows.Close()
	}
}
