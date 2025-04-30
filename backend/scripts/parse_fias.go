package main

import (
	"backend/settings"
	"database/sql"
	"encoding/xml"
	"fmt"
	_ "github.com/lib/pq"
	"os"
	"sync"
)

// СТРУКТУРЫ ДЛЯ XML ФАЙЛОВ С АДРЕСАМИ

type Item struct {
	ID          int    `xml:"ID,attr"`
	ObjectID    string `xml:"OBJECTID,attr"`
	ParentObjID string `xml:"PARENTOBJID,attr"`
	IsActive    string `xml:"ISACTIVE,attr"`
}

type Items struct {
	Items []Item `xml:"ITEM"`
}

type Object struct {
	ID       int    `xml:"ID,attr"`
	ObjectID string `xml:"OBJECTID,attr"`
	Name     string `xml:"NAME,attr"`
	TypeName string `xml:"TYPENAME,attr"`
	FIAS     string `xml:"OBJECTGUID,attr"`
	IsActual string `xml:"ISACTUAL,attr"`
	IsActive string `xml:"ISACTIVE,attr"`
}

type Objects struct {
	Items []Object `xml:"OBJECT"`
}

type House struct {
	ID       int    `xml:"ID,attr"`
	ObjectID string `xml:"OBJECTID,attr"`
	Name     string `xml:"HOUSENUM,attr"`
	TypeID   int    `xml:"HOUSETYPE,attr"`
	FIAS     string `xml:"OBJECTGUID,attr"`
	IsActual string `xml:"ISACTUAL,attr"`
	IsActive string `xml:"ISACTIVE,attr"`
	AddType1 int    `xml:"ADDTYPE1,attr"`
	AddType2 int    `xml:"ADDTYPE2,attr"`
	AddNum1  string `xml:"ADDNUM1,attr"`
	AddNum2  string `xml:"ADDNUM2,attr"`
}

type Houses struct {
	Items []House `xml:"HOUSE"`
}

type HouseType struct {
	ID        int    `xml:"ID,attr"`
	Name      string `xml:"NAME,attr"`
	ShortName string `xml:"SHORTNAME,attr"`
}

type HouseTypes struct {
	Items []HouseType `xml:"HOUSETYPE"`
}

// Скрипт для заполнения таблиц адресами
func main() {
	// ObjectID города Курган
	var streetTypes []HouseType
	mainObjectID := "731553"

	config := settings.Load("../settings.json")

	// Подключаемся к БД
	Link, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DbHost,
		config.DbPort,
		config.DbUser,
		config.DbPass,
		config.DbName))
	if err != nil {
		fmt.Println(err)
		return
	}

	rows, err := Link.Query(`SELECT * FROM "Street_type"`)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer rows.Close()

	for rows.Next() {
		var streetType HouseType
		err = rows.Scan(
			&streetType.ID,
			&streetType.Name,
			&streetType.ShortName,
		)

		streetTypes = append(streetTypes, streetType)
	}

	streetStmt, err := Link.Prepare(`INSERT INTO "Street"(name, type_id, fias_id) VALUES ($1, $2, $3) RETURNING id`)
	if err != nil {
		fmt.Println(err)
		return
	}
	houseStmt, err := Link.Prepare(`INSERT INTO "House"(name, type_id, fias_id, street_id) VALUES ($1, $2, $3, $4) RETURNING id`)
	if err != nil {
		fmt.Println(err)
		return
	}

	insertTypes(Link)

	// Открываем файл с улицами и т.д.
	file, err := os.Open("import/AS_ADDR_OBJ.XML")
	if err != nil {
		fmt.Println("Ошибка открытия файла:", err)
		return
	}
	defer file.Close()

	// Декодируем XML файл и заносим все в переменную "objects"
	var objects Objects
	err = xml.NewDecoder(file).Decode(&objects)
	if err != nil {
		fmt.Println("Ошибка декодирования XML:", err)
		return
	}

	// Открываем файл с домами
	file, err = os.Open("import/AS_HOUSES.XML")
	if err != nil {
		fmt.Println("Ошибка открытия файла:", err)
		return
	}
	defer file.Close()

	// Декодируем XML файл и заносим все в переменную "houses"
	var houses Houses
	err = xml.NewDecoder(file).Decode(&houses)
	if err != nil {
		fmt.Println("Ошибка декодирования XML:", err)
		return
	}

	// Открываем файл со связями объектов
	file, err = os.Open("import/AS_MUN_HIERARCHY.XML")
	if err != nil {
		fmt.Println("Ошибка открытия файла:", err)
		return
	}
	defer file.Close()

	// Декодируем XML файл и заносим все в переменную "items"
	var items Items
	err = xml.NewDecoder(file).Decode(&items)
	if err != nil {
		fmt.Println("Ошибка декодирования XML:", err)
		return
	}

	// Открываем файл с подтипами домов
	file, err = os.Open("import/AS_ADDHOUSE_TYPES.XML")
	if err != nil {
		fmt.Println("Ошибка открытия файла:", err)
		return
	}
	defer file.Close()

	// Декодируем XML файл и заносим все в переменную "addTypes"
	var addTypes HouseTypes
	err = xml.NewDecoder(file).Decode(&addTypes)
	if err != nil {
		fmt.Println("Ошибка декодирования XML:", err)
		return
	}

	fmt.Println(addTypes.Items)

	// Справочник для уже добавленных адресов
	//alreadyExist := make(map[string]bool)
	var alreadyExist sync.Map

	var wg sync.WaitGroup
	var mu sync.Mutex

	// Проходимся циклом по связям для поиска записей с родителем "Курган"
	for _, item0 := range items.Items {
		// Проверям является ли родителем "Курган" и актуальна ли запись
		if item0.ParentObjID == mainObjectID && item0.IsActive == "1" {
			lp := item0

			// Ищем объект из совпавшей связи
			wg.Add(1)
			go func(item0 Item) {
				defer wg.Done()

				//fmt.Println(item0.ObjectID)

				for _, object := range objects.Items {
					// Проверяем не добавили ли мы уже этот объект
					_, ok := alreadyExist.Load(object.FIAS)

					if item0.ObjectID == object.ObjectID && object.IsActive == "1" && object.IsActual == "1" && !ok {

						alreadyExist.Store(object.FIAS, true)

						streetType := 0

						for _, item := range streetTypes {
							if item.ShortName == object.TypeName {
								streetType = item.ID
								break
							}
						}

						var objectID int
						mu.Lock()
						err = streetStmt.QueryRow(object.Name, streetType, object.FIAS).Scan(&objectID)
						mu.Unlock()
						if err != nil {
							fmt.Println("Ошибка выполнения запроса street:", err)
							return
						}

						// Ищем дома у найденого объекта
						for _, item1 := range items.Items {
							if item1.ParentObjID == object.ObjectID && item1.IsActive == "1" {
								for _, house := range houses.Items {

									_, ok = alreadyExist.Load(house.FIAS)

									if item1.ObjectID == house.ObjectID && house.IsActive == "1" && house.IsActual == "1" && !ok {

										alreadyExist.Store(house.FIAS, true)
										//fmt.Println(house.AddType1)

										houseName := house.Name

										if house.AddType1 > 0 {
											for _, addType := range addTypes.Items {
												if addType.ID == house.AddType1 {
													houseName += " " + addType.ShortName + " " + house.AddNum1
												}
											}
										}

										if house.AddType2 > 0 {
											for _, addType := range addTypes.Items {
												if addType.ID == house.AddType2 {
													houseName += " " + addType.ShortName + " " + house.AddNum2
												}
											}
										}

										var houseID int
										mu.Lock()
										err = houseStmt.QueryRow(houseName, house.TypeID, house.FIAS, objectID).Scan(&houseID)
										mu.Unlock()
										if err != nil {
											fmt.Println("Ошибка выполнения запроса house:", err)
											return
										}

										break
									}
								}
							}
						}
						break
					}
				}
			}(lp)
		}
	}

	wg.Wait()
}

// Заносит в БД типы домов и квартир
func insertTypes(Link *sql.DB) {
	// Открываем файл с типами домов
	file, err := os.Open("import/AS_HOUSE_TYPES.XML")
	if err != nil {
		fmt.Println("Ошибка открытия файла:", err)
		return
	}
	defer file.Close()

	// Декодируем XML файл и заносим все в переменную "houseTypes"
	var houseTypes HouseTypes
	err = xml.NewDecoder(file).Decode(&houseTypes)
	if err != nil {
		fmt.Println("Ошибка декодирования XML:", err)
		return
	}

	// Заносим типы домов в БД
	for _, houseType := range houseTypes.Items {
		_, err = Link.Exec(`INSERT INTO "House_type"(name, short_name) VALUES ($1, $2)`, houseType.Name, houseType.ShortName)
		if err != nil {
			fmt.Println("Ошибка выполнения запроса:", err)
			return
		}
	}
}
