package database

import (
	"backend/utils"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type AddressElementType struct {
	ID        int
	Name      string
	ShortName string
}

type AddressElement struct {
	ID   int
	Name string
	Type AddressElementType
	FIAS sql.NullString
}

type Address struct {
	Street         AddressElement
	House          AddressElement
	FileAmount     int
	NodeAmount     int
	HardwareAmount int
}

type Search struct {
	Text   string
	Limit  int
	Offset int
}

var addressElementTypeMap map[string]map[string]bool

func prepareHouse() []string {
	var e error
	errorsList := make([]string, 0)
	addressElementTypeMap = make(map[string]map[string]bool)

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

	query["GET_SUGGESTIONS"], e = Link.Prepare(`
		SELECT s.name, s.type_id, st.short_name, h.id, h.name, h.type_id, ht.short_name, 
		       (SELECT COUNT(*) FROM "House_files" AS f
			                        WHERE f.house_id = h.id ),
				(SELECT COUNT(*) FROM "Node" AS n
										WHERE n.house_id = h.id ),
		    (SELECT COUNT(*) FROM "Hardware" AS hd
		                     JOIN "Node" AS n ON hd.node_id = n.id
										WHERE n.house_id = h.id )
        FROM "Street" AS s
        JOIN "House" AS h ON s.id = h.street_id
        JOIN "Street_type" AS st ON s.type_id = st.id
        JOIN "House_type" AS ht ON h.type_id = ht.id
        WHERE s.name ILIKE '%' || $1 || '%'
          AND (h.name ILIKE '%' || $2 || '%' OR $2 = '')
        ORDER BY 
            CASE 
				WHEN h.name = $2 THEN 0               -- точное полное совпадение (например, "3" == "3")
				WHEN h.name ~ ('^' || $2 || '[^0-9]') THEN 1  -- точное числовое совпадение с префиксом (например, "3а" при "3")
				WHEN h.name ILIKE $2 || '%' THEN 2     -- начинается с номера (например, "3" соответствует "3а")
				ELSE 3                                 -- частичное совпадение
    		END,
			LENGTH(h.name),                           -- сортировка по длине
			h.name
        OFFSET $3
		LIMIT $4
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["GET_SUGGESTIONS_COUNT"], e = Link.Prepare(`
		SELECT COUNT(h.name)
        FROM "Street" AS s
        JOIN "House" AS h ON s.id = h.street_id
        WHERE s.name ILIKE '%' || $1 || '%'
          AND (h.name ILIKE '%' || $2 || '%' OR $2 = '')
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["GET_HOUSE"], e = Link.Prepare(`
		SELECT s.name, s.type_id, st.short_name, h.name, h.type_id, ht.short_name
        FROM "House" AS h
        JOIN "Street" AS s ON s.id = h.street_id
        JOIN "Street_type" AS st ON s.type_id = st.id
        JOIN "House_type" AS ht ON h.type_id = ht.id
        WHERE h.id = $1
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["GET_HOUSES"], e = Link.Prepare(`
		SELECT s.name, s.type_id, st.short_name, h.id, h.name, h.type_id, ht.short_name, 
		        (SELECT COUNT(*) FROM "House_files" AS f
			                        	WHERE f.house_id = h.id ),
				(SELECT COUNT(*) FROM "Node" AS n
										WHERE n.house_id = h.id ),
				(SELECT COUNT(*) FROM "Hardware" AS hd
								 JOIN "Node" AS n ON hd.node_id = n.id
								 WHERE n.house_id = h.id )
        FROM "Street" AS s
        JOIN "House" AS h ON s.id = h.street_id
        JOIN "Street_type" AS st ON s.type_id = st.id
        JOIN "House_type" AS ht ON h.type_id = ht.id
        WHERE EXISTS (
			SELECT 1 
			FROM "House_files" AS f
			WHERE f.house_id = h.id
		) OR EXISTS (
			SELECT 1 
			FROM "Node" AS n
			WHERE n.house_id = h.id
		) OR EXISTS (
			SELECT 1 
			FROM "Hardware" AS hd
			JOIN "Node" AS n ON hd.node_id = n.id
			WHERE n.house_id = h.id
		)
        OFFSET $1
		LIMIT 20
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["GET_HOUSES_COUNT"], e = Link.Prepare(`
		SELECT COUNT(*)
        FROM "Street" AS s
        JOIN "House" AS h ON s.id = h.street_id
        WHERE EXISTS (
			SELECT 1 
			FROM "House_files" AS f
			WHERE f.house_id = h.id
		) OR EXISTS (
			SELECT 1 
			FROM "Node" AS n
			WHERE n.house_id = h.id
		) OR EXISTS (
			SELECT 1 
			FROM "Hardware" AS hd
			JOIN "Node" AS n ON hd.node_id = n.id
			WHERE n.house_id = h.id
		)
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	return errorsList
}

func GetHouses(offset int) ([]Address, int, error) {
	stmt, ok := query["GET_HOUSES"]
	if !ok {
		err := "запрос не подготовлен"
		utils.Logger.Println(err)
		return nil, 0, errors.New(err)
	}

	count, err := countRecord("GET_HOUSES_COUNT", nil)
	if err != nil {
		return nil, 0, err
	}

	rows, err := stmt.Query(offset)
	if err != nil {
		utils.Logger.Println(err)
		return nil, 0, err
	}
	defer rows.Close()

	var addresses []Address

	for rows.Next() {
		var address Address

		err = rows.Scan(
			&address.Street.Name,
			&address.Street.Type.ID,
			&address.Street.Type.ShortName,
			&address.House.ID,
			&address.House.Name,
			&address.House.Type.ID,
			&address.House.Type.ShortName,
			&address.FileAmount,
			&address.NodeAmount,
			&address.HardwareAmount,
		)
		if err != nil {
			utils.Logger.Println(err)
			return nil, 0, err
		}

		addresses = append(addresses, address)
	}

	return addresses, count, nil
}

func (address *Address) GetHouse() error {
	stmt, ok := query["GET_HOUSE"]
	if !ok {
		err := "запрос не подготовлен"
		utils.Logger.Println(err)
		return errors.New(err)
	}

	row := stmt.QueryRow(address.House.ID)

	err := row.Scan(
		&address.Street.Name,
		&address.Street.Type.ID,
		&address.Street.Type.ShortName,
		&address.House.Name,
		&address.House.Type.ID,
		&address.House.Type.ShortName,
	)
	if err != nil {
		utils.Logger.Println(err)
		return err
	}

	return nil
}

func countSuggestions(streetPart, housePart string) (int, error) {
	stmt, ok := query["GET_SUGGESTIONS_COUNT"]
	if !ok {
		err := "запрос не подготовлен"
		utils.Logger.Println(err)
		return 0, errors.New(err)
	}

	row := stmt.QueryRow(streetPart, housePart)

	var count int
	err := row.Scan(&count)
	if err != nil {
		utils.Logger.Println(err)
		return 0, err
	}

	return count, nil
}

func GetSuggestions(search string, offset int, limit int) ([]Address, int, error) {
	streetPart, housePart := parseAddress(search)

	stmt, ok := query["GET_SUGGESTIONS"]
	if !ok {
		err := "запрос GET_SUGGESTIONS не подготовлен"
		utils.Logger.Println(err)
		return nil, 0, errors.New(err)
	}

	var count int
	var err error

	if limit > 10 {
		count, err = countSuggestions(streetPart, housePart)
		if err != nil {
			return nil, 0, err
		}
	}

	rows, err := stmt.Query(streetPart, housePart, offset, limit)
	if err != nil {
		utils.Logger.Println(err)
		return nil, 0, err
	}
	defer rows.Close()

	var addresses []Address
	for rows.Next() {
		var address Address
		err = rows.Scan(
			&address.Street.Name,
			&address.Street.Type.ID,
			&address.Street.Type.ShortName,
			&address.House.ID,
			&address.House.Name,
			&address.House.Type.ID,
			&address.House.Type.ShortName,
			&address.FileAmount,
			&address.NodeAmount,
			&address.HardwareAmount,
		)
		if err != nil {
			utils.Logger.Println(err)
			return nil, 0, err
		}

		addresses = append(addresses, address)
	}

	return addresses, count, nil
}

func parseAddress(input string) (streetPart string, housePart string) {
	// Убираем лишние пробелы и разделяем строку на слова
	cleanedInput := strings.ReplaceAll(input, ",", "")
	words := strings.Fields(cleanedInput)

	streetNameParts := []string{}
	houseDetected := false

	for _, word := range words {
		lowerWord := strings.ToLower(word)

		// Если слово является номером дома
		if !houseDetected && len(streetNameParts) > 0 && isHouseNumber(lowerWord) {
			housePart = lowerWord
			houseDetected = true
			continue
		}

		// Если слово относится к типам улиц или домов, пропускаем его
		if addressElementTypeMap["STREET_TYPE"][lowerWord] || addressElementTypeMap["HOUSE_TYPE"][lowerWord] {
			continue
		}

		// Если дом не был обнаружен, добавляем слово к названию улицы
		if !houseDetected {
			streetNameParts = append(streetNameParts, word)
		}
	}

	streetPart = strings.Join(streetNameParts, " ")
	return streetPart, housePart
}

func isHouseNumber(word string) bool {
	matched, _ := regexp.MatchString(`^\d+[а-яА-Я]?$`, word)
	return matched
}

func countRecord(key string, param interface{}) (int, error) {
	stmt, ok := query[key]
	if !ok {
		err := errors.New("запрос " + key + " не подготовлен")
		utils.Logger.Println(err)
		return 0, err
	}

	var row *sql.Row
	if param != nil {
		row = stmt.QueryRow(param)
	} else {
		row = stmt.QueryRow()
	}

	var count int
	if err := row.Scan(&count); err != nil {
		utils.Logger.Println(err)
		return 0, err
	}

	return count, nil
}

func LoadAddressElementTypeMap(m map[string]map[string]bool) {
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
			var addressElementType AddressElementType
			e = rows.Scan(
				&addressElementType.ID,
				&addressElementType.Name,
				&addressElementType.ShortName,
			)
			if e != nil {
				utils.Logger.Println(e)
				fmt.Println(e)
				return
			}

			m[key][addressElementType.Name] = true
			m[key][addressElementType.ShortName] = true
		}

		_ = rows.Close()
	}
}
