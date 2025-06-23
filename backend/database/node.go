package database

import (
	"backend/models"
	"database/sql"
	"errors"
	"github.com/lib/pq"
)

type NodeRepository interface {
	GetNodesByIDs(nodeIDs []int32) ([]models.Node, error)
	EditNode(node *models.Node) error
	CreateNode(node *models.Node) error
	GetNode(node *models.Node) error
	GetNodes(offset int, onlyActive bool, houseID int) ([]models.Node, int, error)
	ValidateNode(node models.Node) bool
	DeleteNode(nodeID int) error
	GetNodesForIndex() ([]models.Node, error)
}

type DefaultNodeRepository struct {
	Database Database
}

func (r *DefaultNodeRepository) GetNodesForIndex() ([]models.Node, error) {
	stmt, ok := r.Database.GetQuery("GET_NODES_FOR_INDEX")
	if !ok {
		return nil, errors.New("query GET_NODES_FOR_INDEX is not prepare")
	}

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var nodes []models.Node

	for rows.Next() {
		var node models.Node

		if err = rows.Scan(
			&node.ID,
			&node.Name,
			&node.Zone,
			&node.Owner.Value,
			&node.HouseId,
			&node.IsDelete,
			&node.IsPassive,
		); err != nil {
			return nil, err
		}

		nodes = append(nodes, node)
	}

	return nodes, nil
}

func (r *DefaultNodeRepository) DeleteNode(nodeID int) error {
	stmt, ok := r.Database.GetQuery("DELETE_NODE")
	if !ok {
		return errors.New("query DELETE_NODE is not prepare")
	}

	_, err := stmt.Exec(nodeID)
	if err != nil {
		return err
	}

	return nil
}

func (r *DefaultNodeRepository) GetNodesByIDs(nodeIDs []int32) ([]models.Node, error) {
	stmt, ok := r.Database.GetQuery("GET_NODES_BY_IDS")
	if !ok {
		return nil, errors.New("query GET_NODES_BY_IDS is not prepare")
	}

	rows, err := stmt.Query(pq.Array(nodeIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orderedNodes := make([]models.Node, 0, len(nodeIDs))
	nodeMap := make(map[int32]models.Node)

	for rows.Next() {
		var node models.Node

		if err = rows.Scan(
			&node.ID,
			&node.HouseId,
			&node.Owner.ID,
			&node.Name,
			&node.Zone,
			&node.IsPassive,
			&node.Owner.Value,
		); err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}

		nodeMap[int32(node.ID)] = node
	}

	for _, id := range nodeIDs {
		if node, ok := nodeMap[id]; ok {
			orderedNodes = append(orderedNodes, node)
		}
	}

	return orderedNodes, nil
}

func (r *DefaultNodeRepository) EditNode(node *models.Node) error {
	stmt, ok := r.Database.GetQuery("EDIT_NODE")
	if !ok {
		return errors.New("query EDIT_NODE is not prepare")
	}

	var parentID interface{}
	var typeID interface{}

	if node.Parent != nil && !node.IsPassive {
		parentID = node.Parent.ID
	}

	if node.Type != nil && !node.IsPassive {
		typeID = node.Type.ID
	}

	_, err := stmt.Exec(
		node.ID,
		parentID,
		typeID,
		node.Owner.ID,
		node.Name,
		node.Zone,
		node.Placement,
		node.Supply,
		node.Access,
		node.Description,
		node.UpdatedAt,
		node.Address.House.Id,
		node.IsPassive,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *DefaultNodeRepository) CreateNode(node *models.Node) error {
	stmt, ok := r.Database.GetQuery("CREATE_NODE")
	if !ok {
		return errors.New("query CREATE_NODE is not prepare")
	}

	var parentID interface{}
	var typeID interface{}

	if node.Parent != nil && !node.IsPassive {
		parentID = node.Parent.ID
	}

	if node.Type != nil && !node.IsPassive {
		typeID = node.Type.ID
	}

	if err := stmt.QueryRow(
		parentID,
		node.HouseId,
		typeID,
		node.Owner.ID,
		node.Name,
		node.Zone,
		node.Placement,
		node.Supply,
		node.Access,
		node.Description,
		node.CreatedAt,
		node.UpdatedAt,
		node.IsPassive,
	).Scan(&node.ID); err != nil {
		return err
	}

	return nil
}

func (r *DefaultNodeRepository) GetNode(node *models.Node) error {
	stmt, ok := r.Database.GetQuery("GET_NODE")
	if !ok {
		err := errors.New("query GET_NODE is not prepare")
		return err
	}

	var (
		parentID   sql.NullInt64
		parentName sql.NullString
		typeID     sql.NullInt32
		typeValue  sql.NullString
	)

	if err := stmt.QueryRow(node.ID).Scan(
		&node.ID,
		&parentID,
		&node.HouseId,
		&typeID,
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
		&node.IsPassive,
		&typeValue,
		&node.Owner.Value,
		&parentName,
	); err != nil {
		return err
	}

	if parentID.Valid {
		node.Parent = &models.Node{ID: int(parentID.Int64), Name: parentName.String}
	}

	if typeID.Valid {
		node.Type = &models.Reference{ID: int(typeID.Int32), Value: typeValue.String}
	}

	return nil
}

func (r *DefaultNodeRepository) GetNodes(offset int, onlyActive bool, houseID int) ([]models.Node, int, error) {
	stmt, ok := r.Database.GetQuery("GET_NODES")
	if !ok {
		return nil, 0, errors.New("query GET_NODES is not prepare")
	}

	rows, err := stmt.Query(offset, onlyActive, houseID)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var nodes []models.Node
	var count int

	for rows.Next() {
		var node models.Node

		if err = rows.Scan(
			&node.ID,
			&node.HouseId,
			&node.Owner.ID,
			&node.Name,
			&node.Zone,
			&node.IsPassive,
			&node.Owner.Value,
			&count,
		); err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, 0, err
		}

		nodes = append(nodes, node)
	}

	return nodes, count, nil
}

func (r *DefaultNodeRepository) ValidateNode(node models.Node) bool {
	if len(node.Name) == 0 || node.HouseId == 0 || node.Owner.ID == 0 || (!node.IsPassive && node.Type != nil && node.Type.ID == 0) {
		return false
	}

	return true
}
