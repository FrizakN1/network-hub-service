package database

import (
	"backend/models"
	"database/sql"
	"errors"
)

type EventRepository interface {
	CreateEvent(event models.Event) error
	GetEvents(from string, id int) ([]models.Event, int, error)
}

type DefaultEventRepository struct {
	Database Database
	Counter  Counter
}

func (r *DefaultEventRepository) CreateEvent(event models.Event) error {
	stmt, ok := r.Database.GetQuery("CREATE_EVENT")
	if !ok {
		return errors.New("запрос CREATE_EVENT не подготовлен")
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
		event.Address.House.ID,
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

func (r *DefaultEventRepository) GetEvents(from string, id int) ([]models.Event, int, error) {
	key := "GET_EVENTS"

	if from != "" {
		key += "_" + from
	}

	stmt, ok := r.Database.GetQuery(key)
	if !ok {
		return nil, 0, errors.New("запрос " + key + " не подготовлен")
	}

	var rows *sql.Rows
	var err error
	var countParam []interface{} = nil

	if id > 0 {
		rows, err = stmt.Query(id)
		countParam = []interface{}{id}
	} else {
		rows, err = stmt.Query()
	}
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var events []models.Event

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
			&event.Address.House.ID,
			&nodeID,
			&hardwareID,
			&event.UserId,
			&event.Description,
			&event.CreatedAt,
			&event.Address.Street.Name,
			&event.Address.Street.Type.ShortName,
			&event.Address.House.Name,
			&event.Address.House.Type.ShortName,
			&nodeName,
			&hardwareTypeTranslateValue,
		); err != nil {
			return nil, 0, err
		}

		if nodeID.Valid {
			event.Node = &models.Node{ID: int(nodeID.Int64), Name: nodeName.String}
		}

		if hardwareID.Valid {
			event.Hardware = &models.Hardware{ID: int(hardwareID.Int64), Type: models.Reference{TranslateValue: hardwareTypeTranslateValue.String}}
		}

		events = append(events, event)
	}

	count, err := r.Counter.countRecords(key+"_COUNT", countParam)
	if err != nil {
		return nil, 0, err
	}

	return events, count, nil
}
