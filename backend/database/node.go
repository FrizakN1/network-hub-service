package database

import (
	"backend/utils"
	"database/sql"
	"encoding/base64"
	"errors"
	"io/ioutil"
)

type Node struct {
	ID          int
	Parent      *Node
	Address     Address
	Type        Enum
	Owner       Enum
	Name        string
	Zone        sql.NullString
	Placement   sql.NullString
	Supply      sql.NullString
	Access      sql.NullString
	Description sql.NullString
	CreatedAt   int64
	UpdatedAt   sql.NullInt64
}

func prepareNodes() []string {
	var e error
	errorsList := make([]string, 0)

	if query == nil {
		query = make(map[string]*sql.Stmt)
	}

	query["GET_NODES"], e = Link.Prepare(`
		SELECT n.*, s.name, st.short_name, h.name, ht.short_name, nt.name, no.name, (SELECT p.name FROM "Node" AS p WHERE p.id = n.parent_id) AS parent_name
		FROM "Node" AS n 
		JOIN "House" AS h ON n.house_id = h.id
        JOIN "Street" AS s ON s.id = h.street_id
        JOIN "Street_type" AS st ON s.type_id = st.id
        JOIN "House_type" AS ht ON h.type_id = ht.id
		JOIN "Node_type" AS nt ON n.type_id = nt.id
		JOIN "Node_owner" AS no ON n.owner_id = no.id
		ORDER BY id
		OFFSET $1
		LIMIT 20
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["GET_NODES_COUNT"], e = Link.Prepare(`
		SELECT COUNT(*) FROM "Node"
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["GET_HOUSE_NODES"], e = Link.Prepare(`
		SELECT n.*, s.name, st.short_name, h.name, ht.short_name, nt.name, no.name, (SELECT p.name FROM "Node" AS p WHERE p.id = n.parent_id) AS parent_name
		FROM "Node" AS n
		JOIN "House" AS h ON n.house_id = h.id
        JOIN "Street" AS s ON s.id = h.street_id
        JOIN "Street_type" AS st ON s.type_id = st.id
        JOIN "House_type" AS ht ON h.type_id = ht.id
		JOIN "Node_type" AS nt ON n.type_id = nt.id
		JOIN "Node_owner" AS no ON n.owner_id = no.id
		WHERE n.house_id = $1 
		ORDER BY id
		OFFSET $2
		LIMIT 20
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["GET_HOUSE_NODES_COUNT"], e = Link.Prepare(`
		SELECT COUNT(*)
		FROM "Node"
		WHERE house_id = $1 
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["GET_NODE"], e = Link.Prepare(`
		SELECT n.*, s.name, st.short_name, h.name, ht.short_name, nt.name, no.name, (SELECT p.name FROM "Node" AS p WHERE p.id = n.parent_id) AS parent_name
		FROM "Node" AS n 
		JOIN "House" AS h ON n.house_id = h.id
        JOIN "Street" AS s ON s.id = h.street_id
        JOIN "Street_type" AS st ON s.type_id = st.id
        JOIN "House_type" AS ht ON h.type_id = ht.id
		JOIN "Node_type" AS nt ON n.type_id = nt.id
		JOIN "Node_owner" AS no ON n.owner_id = no.id
		WHERE n.id = $1
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["GET_NODE_TYPES"], e = Link.Prepare(`
		SELECT * FROM "Node_type"
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["GET_OWNERS"], e = Link.Prepare(`
		SELECT * FROM "Node_owner"
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["CREATE_OWNER"], e = Link.Prepare(`
		INSERT INTO "Node_owner"(name, created_at) VALUES ($1, $2)
		RETURNING id
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["EDIT_OWNER"], e = Link.Prepare(`
		UPDATE "Node_owner" SET name = $2 WHERE id = $1
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["CREATE_NODE_TYPE"], e = Link.Prepare(`
		INSERT INTO "Node_type"(name, created_at) VALUES ($1, $2)
		RETURNING id
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["EDIT_NODE_TYPE"], e = Link.Prepare(`
		UPDATE "Node_type" SET name = $2 WHERE id = $1
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["CREATE_HARDWARE_TYPE"], e = Link.Prepare(`
		INSERT INTO "Hardware_type"(value, translate_value, created_at) VALUES ($1, $2, $3)
		RETURNING id
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["EDIT_HARDWARE_TYPE"], e = Link.Prepare(`
		UPDATE "Hardware_type" SET value = $2, translate_value = $3 WHERE id = $1
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["CREATE_OPERATION_MODE"], e = Link.Prepare(`
		INSERT INTO "Operation_mode"(value, translate_value, created_at) VALUES ($1, $2, $3)
		RETURNING id
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["EDIT_OPERATION_MODE"], e = Link.Prepare(`
		UPDATE "Operation_mode" SET value = $2, translate_value = $3 WHERE id = $1
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["CREATE_NODE"], e = Link.Prepare(`
		INSERT INTO "Node"(parent_id, house_id, type_id, owner_id, name, zone, placement, supply, access, description, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["EDIT_NODE"], e = Link.Prepare(`
		UPDATE "Node" SET parent_id = $2, type_id = $3, owner_id = $4, name = $5, zone = $6, placement = $7, supply = $8,
		                  access = $9, description = $10, updated_at = $11, house_id = $12
		WHERE id = $1	
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["GET_NODE_FILES"], e = Link.Prepare(`
		SELECT * FROM "Node_files" WHERE node_id = $1 AND is_preview_image = $2
		ORDER BY upload_at DESC 
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["GET_SEARCH_NODES"], e = Link.Prepare(`
		SELECT n.*, s.name, st.short_name, h.name, ht.short_name, nt.name, no.name, (SELECT p.name FROM "Node" AS p WHERE p.id = n.parent_id) AS parent_name
		FROM "Node" AS n 
		JOIN "House" AS h ON n.house_id = h.id
        JOIN "Street" AS s ON s.id = h.street_id
        JOIN "Street_type" AS st ON s.type_id = st.id
        JOIN "House_type" AS ht ON h.type_id = ht.id
		JOIN "Node_type" AS nt ON n.type_id = nt.id
		JOIN "Node_owner" AS no ON n.owner_id = no.id
		WHERE n.name ILIKE '%' || $1 || '%'
			OR nt.name ILIKE '%' || $1 || '%'
			OR no.name ILIKE '%' || $1 || '%'
			OR n.zone ILIKE '%' || $1 || '%'
			OR s.name ILIKE '%' || $1 || '%'
			OR st.short_name ILIKE '%' || $1 || '%'
			OR h.name ILIKE '%' || $1 || '%'
			OR ht.short_name ILIKE '%' || $1 || '%'
			OR (st.short_name || ' ' || s.name || ', ' || ht.short_name || ' ' || h.name) ILIKE '%' || $1 || '%'
		ORDER BY id
		OFFSET $2
		LIMIT 20
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["GET_SEARCH_NODES_COUNT"], e = Link.Prepare(`
		SELECT COUNT(n.*)
		FROM "Node" AS n 
		JOIN "House" AS h ON n.house_id = h.id
        JOIN "Street" AS s ON s.id = h.street_id
        JOIN "Street_type" AS st ON s.type_id = st.id
        JOIN "House_type" AS ht ON h.type_id = ht.id
		JOIN "Node_type" AS nt ON n.type_id = nt.id
		JOIN "Node_owner" AS no ON n.owner_id = no.id
		WHERE n.name ILIKE '%' || $1 || '%'
			OR nt.name ILIKE '%' || $1 || '%'
			OR no.name ILIKE '%' || $1 || '%'
			OR n.zone ILIKE '%' || $1 || '%'
			OR s.name ILIKE '%' || $1 || '%'
			OR st.short_name ILIKE '%' || $1 || '%'
			OR h.name ILIKE '%' || $1 || '%'
			OR ht.short_name ILIKE '%' || $1 || '%'
			OR (st.short_name || ' ' || s.name || ', ' || ht.short_name || ' ' || h.name) ILIKE '%' || $1 || '%'
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	return errorsList
}

func GetSearchNodes(search string, offset int) ([]Node, int, error) {
	stmt, ok := query["GET_SEARCH_NODES"]
	if !ok {
		err := errors.New("запрос GET_SEARCH_NODES не подготовлен")
		utils.Logger.Println(err)
		return nil, 0, err
	}

	count, err := countRecord("GET_SEARCH_NODES_COUNT", search)
	if err != nil {
		utils.Logger.Println(err)
		return nil, 0, err
	}

	rows, err := stmt.Query(search, offset)
	if err != nil {
		utils.Logger.Println(err)
		return nil, 0, err
	}
	defer rows.Close()

	var nodes []Node

	for rows.Next() {
		var (
			node       Node
			parentID   sql.NullInt64
			parentName sql.NullString
		)

		if err = rows.Scan(
			&node.ID,
			&parentID,
			&node.Address.House.ID,
			&node.Type.ID,
			&node.Owner.ID,
			&node.Name,
			&node.Zone,
			&node.Placement,
			&node.Supply,
			&node.Access,
			&node.Description,
			&node.CreatedAt,
			&node.UpdatedAt,
			&node.Address.Street.Name,
			&node.Address.Street.Type.ShortName,
			&node.Address.House.Name,
			&node.Address.House.Type.ShortName,
			&node.Type.Name,
			&node.Owner.Name,
			&parentName,
		); err != nil && !errors.Is(err, sql.ErrNoRows) {
			utils.Logger.Println(err)
			return nil, 0, err
		}

		if parentID.Valid {
			node.Parent = &Node{ID: int(parentID.Int64), Name: parentName.String}
		}

		nodes = append(nodes, node)
	}

	return nodes, count, nil
}

func GetNodeFiles(nodeID int, onlyImage bool) ([]File, error) {
	stmt, ok := query["GET_NODE_FILES"]
	if !ok {
		err := errors.New("запрос GET_NODE_FILES не подготовлен")
		utils.Logger.Println(err)
		return nil, err
	}

	rows, err := stmt.Query(nodeID, onlyImage)
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
			&file.Node.ID,
			&file.Path,
			&file.Name,
			&file.UploadAt,
			&file.InArchive,
			&file.IsPreviewImage,
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

func (node *Node) EditNode() error {
	stmt, ok := query["EDIT_NODE"]
	if !ok {
		err := errors.New("запрос EDIT_NODE не подготовлен")
		utils.Logger.Println(err)
		return err
	}

	var parentID interface{}

	if node.Parent != nil {
		parentID = node.Parent.ID
	}

	_, err := stmt.Exec(
		node.ID,
		parentID,
		node.Type.ID,
		node.Owner.ID,
		node.Name,
		node.Zone,
		node.Placement,
		node.Supply,
		node.Access,
		node.Description,
		node.UpdatedAt,
		node.Address.House.ID,
	)
	if err != nil {
		utils.Logger.Println(err)
		return err
	}

	return nil
}

func (node *Node) CreateNode() error {
	stmt, ok := query["CREATE_NODE"]
	if !ok {
		err := errors.New("запрос CREATE_NODE не подготовлен")
		utils.Logger.Println(err)
		return err
	}

	var parentID interface{}

	if node.Parent != nil {
		parentID = node.Parent.ID
	}

	if err := stmt.QueryRow(
		parentID,
		node.Address.House.ID,
		node.Type.ID,
		node.Owner.ID,
		node.Name,
		node.Zone,
		node.Placement,
		node.Supply,
		node.Access,
		node.Description,
		node.CreatedAt,
		node.UpdatedAt,
	).Scan(&node.ID); err != nil {
		utils.Logger.Println(err)
		return err
	}

	return nil
}

func (referenceRecord *Enum) EditReferenceRecord(reference string) error {
	stmt, ok := query["EDIT_"+reference]
	if !ok {
		err := errors.New("запрос EDIT_" + reference + " не подготовлен")
		utils.Logger.Println(err)
		return err
	}

	switch reference {
	case "NODE_TYPE":
	case "OWNER":
		_, err := stmt.Exec(referenceRecord.ID, referenceRecord.Name)
		if err != nil {
			utils.Logger.Println(err)
			return err
		}
		break
	case "HARDWARE_TYPE":
	case "OPERATION_MODE":
		_, err := stmt.Exec(referenceRecord.ID, referenceRecord.Value, referenceRecord.TranslateValue)
		if err != nil {
			utils.Logger.Println(err)
			return err
		}
		break
	}

	return nil
}

func (referenceRecord *Enum) CreateReferenceRecord(reference string) error {
	stmt, ok := query["CREATE_"+reference]
	if !ok {
		err := errors.New("запрос CREATE_" + reference + " не подготовлен")
		utils.Logger.Println(err)
		return err
	}

	if reference == "NODE_TYPE" || reference == "OWNER" {
		if err := stmt.QueryRow(
			referenceRecord.Name,
			referenceRecord.CreatedAt,
		).Scan(&referenceRecord.ID); err != nil {
			utils.Logger.Println(err)
			return err
		}
	} else {
		if err := stmt.QueryRow(
			referenceRecord.Value,
			referenceRecord.TranslateValue,
			referenceRecord.CreatedAt,
		).Scan(&referenceRecord.ID); err != nil {
			utils.Logger.Println(err)
			return err
		}
	}

	return nil
}

func GetNodeEnums(key string) ([]Enum, error) {
	stmt, ok := query["GET_"+key]
	if !ok {
		err := errors.New("запрос GET_" + key + " не подготовлен")
		utils.Logger.Println(err)
		return nil, err
	}

	rows, err := stmt.Query()
	if err != nil {
		utils.Logger.Println(err)
		return nil, err
	}
	defer rows.Close()

	var enums []Enum

	for rows.Next() {
		var enum Enum

		if err = rows.Scan(
			&enum.ID,
			&enum.Name,
			&enum.CreatedAt,
		); err != nil {
			utils.Logger.Println(err)
			return nil, err
		}

		enums = append(enums, enum)
	}

	return enums, nil
}

func (node *Node) GetNode() error {
	stmt, ok := query["GET_NODE"]
	if !ok {
		err := errors.New("запрос GET_NODE не подготовлен")
		utils.Logger.Println(err)
		return err
	}

	var (
		parentID   sql.NullInt64
		parentName sql.NullString
	)

	if err := stmt.QueryRow(node.ID).Scan(
		&node.ID,
		&parentID,
		&node.Address.House.ID,
		&node.Type.ID,
		&node.Owner.ID,
		&node.Name,
		&node.Zone,
		&node.Placement,
		&node.Supply,
		&node.Access,
		&node.Description,
		&node.CreatedAt,
		&node.UpdatedAt,
		&node.Address.Street.Name,
		&node.Address.Street.Type.ShortName,
		&node.Address.House.Name,
		&node.Address.House.Type.ShortName,
		&node.Type.Name,
		&node.Owner.Name,
		&parentName,
	); err != nil {
		utils.Logger.Println(err)
		return err
	}

	if parentID.Valid {
		node.Parent = &Node{ID: int(parentID.Int64), Name: parentName.String}
	}

	return nil
}

func GetHouseNodes(houseID int, offset int) ([]Node, int, error) {
	stmt, ok := query["GET_HOUSE_NODES"]
	if !ok {
		err := errors.New("запрос GET_HOUSE_NODES не подготовлен")
		utils.Logger.Println(err)
		return nil, 0, err
	}

	count, err := countRecord("GET_HOUSE_NODES_COUNT", houseID)
	if err != nil {
		utils.Logger.Println(err)
		return nil, 0, err
	}

	rows, err := stmt.Query(houseID, offset)
	if err != nil {
		utils.Logger.Println(err)
		return nil, 0, err
	}
	defer rows.Close()

	var nodes []Node

	for rows.Next() {
		var (
			node       Node
			parentID   sql.NullInt64
			parentName sql.NullString
		)

		if err = rows.Scan(
			&node.ID,
			&parentID,
			&node.Address.House.ID,
			&node.Type.ID,
			&node.Owner.ID,
			&node.Name,
			&node.Zone,
			&node.Placement,
			&node.Supply,
			&node.Access,
			&node.Description,
			&node.CreatedAt,
			&node.UpdatedAt,
			&node.Address.Street.Name,
			&node.Address.Street.Type.ShortName,
			&node.Address.House.Name,
			&node.Address.House.Type.ShortName,
			&node.Type.Name,
			&node.Owner.Name,
			&parentName,
		); err != nil && !errors.Is(err, sql.ErrNoRows) {
			utils.Logger.Println(err)
			return nil, 0, err
		}

		if parentID.Valid {
			node.Parent = &Node{ID: int(parentID.Int64), Name: parentName.String}
		}

		nodes = append(nodes, node)
	}

	return nodes, count, nil
}

func GetNodes(offset int) ([]Node, int, error) {
	stmt, ok := query["GET_NODES"]
	if !ok {
		err := errors.New("запрос GET_NODES не подготовлен")
		utils.Logger.Println(err)
		return nil, 0, err
	}

	count, err := countRecord("GET_NODES_COUNT", nil)
	if err != nil {
		utils.Logger.Println(err)
		return nil, 0, err
	}

	rows, err := stmt.Query(offset)
	if err != nil {
		utils.Logger.Println(err)
		return nil, 0, err
	}
	defer rows.Close()

	var nodes []Node

	for rows.Next() {
		var (
			node       Node
			parentID   sql.NullInt64
			parentName sql.NullString
		)

		if err = rows.Scan(
			&node.ID,
			&parentID,
			&node.Address.House.ID,
			&node.Type.ID,
			&node.Owner.ID,
			&node.Name,
			&node.Zone,
			&node.Placement,
			&node.Supply,
			&node.Access,
			&node.Description,
			&node.CreatedAt,
			&node.UpdatedAt,
			&node.Address.Street.Name,
			&node.Address.Street.Type.ShortName,
			&node.Address.House.Name,
			&node.Address.House.Type.ShortName,
			&node.Type.Name,
			&node.Owner.Name,
			&parentName,
		); err != nil && !errors.Is(err, sql.ErrNoRows) {
			utils.Logger.Println(err)
			return nil, 0, err
		}

		if parentID.Valid {
			node.Parent = &Node{ID: int(parentID.Int64), Name: parentName.String}
		}

		nodes = append(nodes, node)
	}

	return nodes, count, nil
}

func (node *Node) ValidateNode() bool {
	if len(node.Name) == 0 || node.Address.House.ID == 0 || node.Type.ID == 0 || node.Owner.ID == 0 {
		return false
	}

	return true
}
