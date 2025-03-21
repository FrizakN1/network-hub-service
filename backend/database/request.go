package database

import (
	"backend/utils"
	"database/sql"
	"encoding/base64"
	"errors"
	"io/ioutil"
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
	Street     AddressElement
	House      AddressElement
	FileAmount int
}

type File struct {
	ID        int
	House     AddressElement
	Path      string
	Name      string
	UploadAt  int64
	Data      string
	InArchive bool
}

type Search struct {
	Text   string
	Limit  int
	Offset int
}

func prepareRequests() []string {
	var e error
	errorsList := make([]string, 0)

	if query == nil {
		query = make(map[string]*sql.Stmt)
	}

	query["GET_SUGGESTIONS"], e = Link.Prepare(`
		SELECT s.name, s.type_id, st.short_name, h.id, h.name, h.type_id, ht.short_name, 
		       (SELECT COUNT(*) FROM "Files" AS f
			                        WHERE f.house_id = h.id ) 
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

	query["CREATE_FILE"], e = Link.Prepare(`
		INSERT INTO "Files"(house_id, file_path, file_name, upload_at, in_archive) 
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["GET_FILES"], e = Link.Prepare(`
		SELECT * FROM "Files" WHERE house_id = $1
		ORDER BY upload_at DESC 
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["ARCHIVE_FILE"], e = Link.Prepare(`
		UPDATE "Files" SET in_archive = $2 WHERE id = $1
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["DELETE_FILE"], e = Link.Prepare(`
		DELETE FROM "Files" WHERE id = $1
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["GET_LIST"], e = Link.Prepare(`
		SELECT s.name, s.type_id, st.short_name, h.id, h.name, h.type_id, ht.short_name, 
		       (SELECT COUNT(*) FROM "Files" AS f
			                        WHERE f.house_id = h.id ) 
        FROM "Street" AS s
        JOIN "House" AS h ON s.id = h.street_id
        JOIN "Street_type" AS st ON s.type_id = st.id
        JOIN "House_type" AS ht ON h.type_id = ht.id
        WHERE EXISTS (
			SELECT 1 
			FROM "Files" AS f
			WHERE f.house_id = h.id
		)
        OFFSET $1
		LIMIT 20
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["GET_LIST_COUNT"], e = Link.Prepare(`
		SELECT COUNT(*)
        FROM "Street" AS s
        JOIN "House" AS h ON s.id = h.street_id
        WHERE EXISTS (
			SELECT 1 
			FROM "Files" AS f
			WHERE f.house_id = h.id
		)
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	return errorsList
}

func countList() (int, error) {
	stmt, ok := query["GET_LIST_COUNT"]
	if !ok {
		err := "запрос не подготовлен"
		utils.Logger.Println(err)
		return 0, errors.New(err)
	}

	row := stmt.QueryRow()

	var count int
	err := row.Scan(&count)
	if err != nil {
		utils.Logger.Println(err)
		return 0, err
	}

	return count, nil
}

func GetList(offset int) ([]Address, int, error) {
	stmt, ok := query["GET_LIST"]
	if !ok {
		err := "запрос не подготовлен"
		utils.Logger.Println(err)
		return nil, 0, errors.New(err)
	}

	count, err := countList()
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
		)
		if err != nil {
			utils.Logger.Println(err)
			return nil, 0, err
		}

		addresses = append(addresses, address)
	}

	return addresses, count, nil
}

func (file *File) Delete() error {
	stmt, ok := query["DELETE_FILE"]
	if !ok {
		err := "запрос не подготовлен"
		utils.Logger.Println(err)
		return errors.New(err)
	}

	_, err := stmt.Exec(file.ID)
	if err != nil {
		utils.Logger.Println(err)
		return err
	}

	return nil
}

func (file *File) Archive() error {
	stmt, ok := query["ARCHIVE_FILE"]
	if !ok {
		err := "запрос не подготовлен"
		utils.Logger.Println(err)
		return errors.New(err)
	}

	file.InArchive = !file.InArchive

	_, err := stmt.Exec(file.ID, file.InArchive)
	if err != nil {
		utils.Logger.Println(err)
		return err
	}

	return nil
}

func GetFiles(houseID int) ([]File, error) {
	stmt, ok := query["GET_FILES"]
	if !ok {
		err := "запрос не подготовлен"
		utils.Logger.Println(err)
		return nil, errors.New(err)
	}

	rows, err := stmt.Query(houseID)
	if err != nil {
		utils.Logger.Println(err)
		return nil, err
	}
	defer rows.Close()

	var files []File
	for rows.Next() {
		var file File

		err = rows.Scan(
			&file.ID,
			&file.House.ID,
			&file.Path,
			&file.Name,
			&file.UploadAt,
			&file.InArchive,
		)
		if err != nil {
			utils.Logger.Println(err)
			return nil, err
		}

		var fileData []byte

		fileData, err = ioutil.ReadFile(file.Path)
		if err != nil {
			utils.Logger.Println(err)
			return nil, err
		}

		file.Data = base64.StdEncoding.EncodeToString(fileData)

		files = append(files, file)
	}

	return files, nil
}

func (file *File) Create() error {
	stmt, ok := query["CREATE_FILE"]
	if !ok {
		err := "запрос не подготовлен"
		utils.Logger.Println(err)
		return errors.New(err)
	}

	err := stmt.QueryRow(
		file.House.ID,
		file.Path,
		file.Name,
		file.UploadAt,
		file.InArchive,
	).Scan(&file.ID)
	if err != nil {
		utils.Logger.Println(err)
		return err
	}

	return nil
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

func GetSuggestions(search Search) ([]Address, int, error) {
	streetPart, housePart := parseAddress(search.Text)

	stmt, ok := query["GET_SUGGESTIONS"]
	if !ok {
		err := "запрос не подготовлен"
		utils.Logger.Println(err)
		return nil, 0, errors.New(err)
	}

	var count int
	var err error

	if search.Limit > 10 {
		count, err = countSuggestions(streetPart, housePart)
		if err != nil {
			return nil, 0, err
		}
	}

	rows, err := stmt.Query(streetPart, housePart, search.Offset, search.Limit)
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
		if enumsMap["STREET_TYPE"][lowerWord] || enumsMap["HOUSE_TYPE"][lowerWord] {
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
