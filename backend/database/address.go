package database

import (
	"backend/models"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type AddressRepository interface {
	GetHouses(offset int) ([]models.Address, int, error)
	GetHouse(addresr *models.Address) error
	GetSuggestions(search string, offset int, limit int) ([]models.Address, int, error)
	SetHouseParams(address *models.Address) error
}

type DefaultAddressRepository struct {
	addressElementTypeMap map[string]map[string]struct{}
	Database              Database
}

func (r *DefaultAddressRepository) SetHouseParams(address *models.Address) error {
	stmt, ok := r.Database.GetQuery("SET_HOUSE_PARAMS")
	if !ok {
		return errors.New("query SET_HOUSE_PARAMS is not prepare")
	}

	_, err := stmt.Exec(address.House.ID, address.RoofType.ID, address.WiringType.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *DefaultAddressRepository) GetHouses(offset int) ([]models.Address, int, error) {
	stmt, ok := r.Database.GetQuery("GET_HOUSES")
	if !ok {
		return nil, 0, errors.New("query GET_HOUSES is not prepare")
	}

	rows, err := stmt.Query(offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var addresses []models.Address
	var count int

	for rows.Next() {
		var address models.Address

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
			&count,
		)
		if err != nil {
			return nil, 0, err
		}

		addresses = append(addresses, address)
	}

	return addresses, count, nil
}

func (r *DefaultAddressRepository) GetHouse(address *models.Address) error {
	stmt, ok := r.Database.GetQuery("GET_HOUSE")
	if !ok {
		return errors.New("query GET_HOUSE is not prepare")
	}

	row := stmt.QueryRow(address.House.ID)

	var roofID sql.NullInt16
	var roofValue sql.NullString
	var wiringID sql.NullInt16
	var wiringValue sql.NullString

	err := row.Scan(
		&address.Street.Name,
		&address.Street.Type.ID,
		&address.Street.Type.ShortName,
		&address.House.Name,
		&address.House.Type.ID,
		&address.House.Type.ShortName,
		&roofID,
		&roofValue,
		&wiringID,
		&wiringValue,
	)
	if err != nil {
		return err
	}

	if roofID.Valid {
		address.RoofType = models.Reference{
			ID:    int(roofID.Int16),
			Value: roofValue.String,
		}
	}

	if wiringID.Valid {
		address.WiringType = models.Reference{
			ID:    int(wiringID.Int16),
			Value: wiringValue.String,
		}
	}

	return nil
}

func (r *DefaultAddressRepository) GetSuggestions(search string, offset int, limit int) ([]models.Address, int, error) {
	streetPart, housePart := parseAddress(search, r.addressElementTypeMap)

	stmt, ok := r.Database.GetQuery("GET_SUGGESTIONS")
	if !ok {
		return nil, 0, errors.New("query GET_SUGGESTIONS is not prepare")
	}

	var err error

	rows, err := stmt.Query(streetPart, housePart, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var addresses []models.Address
	var count int

	for rows.Next() {
		var address models.Address
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
			&count,
		)
		if err != nil {
			return nil, 0, err
		}

		addresses = append(addresses, address)
	}

	return addresses, count, nil
}

func parseAddress(input string, addressElementTypeMap map[string]map[string]struct{}) (streetPart string, housePart string) {
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

		_, streetTypeExist := addressElementTypeMap["STREET_TYPE"][lowerWord]
		_, houseTypeExist := addressElementTypeMap["HOUSE_TYPE"][lowerWord]

		// Если слово относится к типам улиц или домов, пропускаем его
		if streetTypeExist || houseTypeExist {
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

func (r *DefaultAddressRepository) LoadAddressElementTypeMap() {
	r.addressElementTypeMap = make(map[string]map[string]struct{}, 0)

	keys := [2]string{
		"STREET_TYPE",
		"HOUSE_TYPE",
	}

	for _, key := range keys {
		stmt, ok := r.Database.GetQuery(key)
		if !ok {
			fmt.Println("ошибка загрузки перечислений, необходимо остановить работу сервера и обратиться к разработчику")
			return
		}

		rows, e := stmt.Query()
		if e != nil {
			fmt.Println(e)
			fmt.Println("ошибка загрузки перечислений, необходимо остановить работу сервера и обратиться к разработчику")
			return
		}

		if r.addressElementTypeMap[key] == nil {
			r.addressElementTypeMap[key] = make(map[string]struct{})
		}

		for rows.Next() {
			var addressElementType models.AddressElementType
			e = rows.Scan(
				&addressElementType.ID,
				&addressElementType.Name,
				&addressElementType.ShortName,
			)
			if e != nil {
				fmt.Println(e)
				return
			}

			r.addressElementTypeMap[key][addressElementType.Name] = struct{}{}
			r.addressElementTypeMap[key][addressElementType.ShortName] = struct{}{}
		}

		_ = rows.Close()
	}
}
