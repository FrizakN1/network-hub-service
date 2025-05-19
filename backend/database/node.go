package database

import (
	"backend/models"
	"database/sql"
	"errors"
)

type NodeRepository interface {
	GetSearchNodes(search string, offset int) ([]models.Node, int, error)
	EditNode(node *models.Node) error
	CreateNode(node *models.Node) error
	GetNode(node *models.Node) error
	GetHouseNodes(houseID int, offset int) ([]models.Node, int, error)
	GetNodes(offset int) ([]models.Node, int, error)
	ValidateNode(node models.Node) bool
	DeleteNode(nodeID int) error
}

type DefaultNodeRepository struct {
	Database Database
	Counter  Counter
}

func (r *DefaultNodeRepository) DeleteNode(nodeID int) error {
	stmt, ok := r.Database.GetQuery("DELETE_NODE")
	if !ok {
		return errors.New("запрос DELETE_NODE не подготовлен")
	}

	_, err := stmt.Exec(nodeID)
	if err != nil {
		return err
	}

	return nil
}

func (r *DefaultNodeRepository) GetSearchNodes(search string, offset int) ([]models.Node, int, error) {
	stmt, ok := r.Database.GetQuery("GET_SEARCH_NODES")
	if !ok {
		return nil, 0, errors.New("запрос GET_SEARCH_NODES не подготовлен")
	}

	count, err := r.Counter.countRecords("GET_SEARCH_NODES_COUNT", []interface{}{search})
	if err != nil {
		return nil, 0, err
	}

	rows, err := stmt.Query(search, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var nodes []models.Node

	for rows.Next() {
		var (
			node       models.Node
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
			&node.IsDelete,
			&node.Address.Street.Name,
			&node.Address.Street.Type.ShortName,
			&node.Address.House.Name,
			&node.Address.House.Type.ShortName,
			&node.Type.Name,
			&node.Owner.Name,
			&parentName,
		); err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, 0, err
		}

		if parentID.Valid {
			node.Parent = &models.Node{ID: int(parentID.Int64), Name: parentName.String}
		}

		nodes = append(nodes, node)
	}

	return nodes, count, nil
}

func (r *DefaultNodeRepository) EditNode(node *models.Node) error {
	stmt, ok := r.Database.GetQuery("EDIT_NODE")
	if !ok {
		return errors.New("запрос EDIT_NODE не подготовлен")
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
		return err
	}

	return nil
}

func (r *DefaultNodeRepository) CreateNode(node *models.Node) error {
	stmt, ok := r.Database.GetQuery("CREATE_NODE")
	if !ok {
		return errors.New("запрос CREATE_NODE не подготовлен")
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
		return err
	}

	return nil
}

func (r *DefaultNodeRepository) GetNode(node *models.Node) error {
	stmt, ok := r.Database.GetQuery("GET_NODE")
	if !ok {
		err := errors.New("запрос GET_NODE не подготовлен")
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
		&node.IsDelete,
		&node.Address.Street.Name,
		&node.Address.Street.Type.ShortName,
		&node.Address.House.Name,
		&node.Address.House.Type.ShortName,
		&node.Type.Name,
		&node.Owner.Name,
		&parentName,
	); err != nil {
		return err
	}

	if parentID.Valid {
		node.Parent = &models.Node{ID: int(parentID.Int64), Name: parentName.String}
	}

	return nil
}

func (r *DefaultNodeRepository) GetHouseNodes(houseID int, offset int) ([]models.Node, int, error) {
	stmt, ok := r.Database.GetQuery("GET_HOUSE_NODES")
	if !ok {
		return nil, 0, errors.New("запрос GET_HOUSE_NODES не подготовлен")
	}

	count, err := r.Counter.countRecords("GET_HOUSE_NODES_COUNT", []interface{}{houseID})
	if err != nil {
		return nil, 0, err
	}

	rows, err := stmt.Query(houseID, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var nodes []models.Node

	for rows.Next() {
		var (
			node       models.Node
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
			&node.IsDelete,
			&node.Address.Street.Name,
			&node.Address.Street.Type.ShortName,
			&node.Address.House.Name,
			&node.Address.House.Type.ShortName,
			&node.Type.Name,
			&node.Owner.Name,
			&parentName,
		); err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, 0, err
		}

		if parentID.Valid {
			node.Parent = &models.Node{ID: int(parentID.Int64), Name: parentName.String}
		}

		nodes = append(nodes, node)
	}

	return nodes, count, nil
}

func (r *DefaultNodeRepository) GetNodes(offset int) ([]models.Node, int, error) {
	stmt, ok := r.Database.GetQuery("GET_NODES")
	if !ok {
		return nil, 0, errors.New("запрос GET_NODES не подготовлен")
	}

	count, err := r.Counter.countRecords("GET_NODES_COUNT", nil)
	if err != nil {
		return nil, 0, err
	}

	rows, err := stmt.Query(offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var nodes []models.Node

	for rows.Next() {
		var (
			node       models.Node
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
			&node.IsDelete,
			&node.Address.Street.Name,
			&node.Address.Street.Type.ShortName,
			&node.Address.House.Name,
			&node.Address.House.Type.ShortName,
			&node.Type.Name,
			&node.Owner.Name,
			&parentName,
		); err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, 0, err
		}

		if parentID.Valid {
			node.Parent = &models.Node{ID: int(parentID.Int64), Name: parentName.String}
		}

		nodes = append(nodes, node)
	}

	return nodes, count, nil
}

func (r *DefaultNodeRepository) ValidateNode(node models.Node) bool {
	if len(node.Name) == 0 || node.Address.House.ID == 0 || node.Type.ID == 0 || node.Owner.ID == 0 {
		return false
	}

	return true
}
