package database

import (
	"backend/models"
	"database/sql"
	"errors"
)

type EventRepository interface {
	CreateEvent(event models.Event) error
	GetEvents(offset int, from string, id int) ([]models.Event, int, error)
}

type DefaultEventRepository struct {
	Database Database
}

func (r *DefaultEventRepository) CreateEvent(event models.Event) error {
	stmt, ok := r.Database.GetQuery("CREATE_EVENT")
	if !ok {
		return errors.New("query CREATE_EVENT is not prepare")
	}

	var nodeID interface{}
	var hardwareID interface{}

	if event.Node != nil {
		nodeID = event.Node.ID
	}

	if event.Hardware != nil {
		hardwareID = event.Hardware.ID
	}

	_, err := stmt.Exec(
		event.HouseId,
		nodeID,
		hardwareID,
		event.UserId,
		event.Description,
		event.CreatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *DefaultEventRepository) GetEvents(offset int, from string, id int) ([]models.Event, int, error) {
	key := "GET_EVENTS"

	if from != "" {
		key += "_" + from
	}

	stmt, ok := r.Database.GetQuery(key)
	if !ok {
		return nil, 0, errors.New("query " + key + " is not prepare")
	}

	var rows *sql.Rows
	var err error

	if id > 0 {
		rows, err = stmt.Query(offset, id)
	} else {
		rows, err = stmt.Query(offset)
	}
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var events []models.Event
	var count int

	for rows.Next() {
		var (
			event                      models.Event
			nodeID                     sql.NullInt64
			hardwareID                 sql.NullInt64
			nodeName                   sql.NullString
			hardwareTypeTranslateValue sql.NullString
		)

		if err = rows.Scan(
			&event.ID,
			&event.HouseId,
			&nodeID,
			&hardwareID,
			&event.UserId,
			&event.Description,
			&event.CreatedAt,
			&nodeName,
			&hardwareTypeTranslateValue,
			&count,
		); err != nil {
			return nil, 0, err
		}

		if nodeID.Valid {
			event.Node = &models.Node{ID: int(nodeID.Int64), Name: nodeName.String}
		}

		if hardwareID.Valid {
			event.Hardware = &models.Hardware{ID: int(hardwareID.Int64), Type: models.Reference{Value: hardwareTypeTranslateValue.String}}
		}

		events = append(events, event)
	}

	return events, count, nil
}
