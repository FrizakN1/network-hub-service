package database

import (
	"backend/utils"
	"database/sql"
	"errors"
)

type Event struct {
	ID          int64
	Address     Address
	Node        *Node
	Hardware    *Hardware
	User        User
	Description string
	CreatedAt   int64
}

func prepareEvent() []string {
	var e error
	errorsList := make([]string, 0)

	if query == nil {
		query = make(map[string]*sql.Stmt)
	}

	query["GET_EVENTS"], e = Link.Prepare(`
		SELECT e.*, u.login, u.name, s.name, st.short_name, h.name, ht.short_name, n.name, hwt.translate_value
		FROM "Event" AS e
		JOIN "House" AS h ON e.house_id = h.id
        JOIN "Street" AS s ON s.id = h.street_id
        JOIN "Street_type" AS st ON s.type_id = st.id
        JOIN "House_type" AS ht ON h.type_id = ht.id
		JOIN "User" AS u ON e.user_id = u.id
		LEFT JOIN "Node" AS n ON e.node_id = n.id
		LEFT JOIN "Hardware" AS hw ON e.hardware_id = hw.id
		LEFT JOIN "Hardware_type" AS hwt ON hw.type_id = hwt.id
		ORDER BY e.created_at DESC
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["GET_EVENTS_COUNT"], e = Link.Prepare(`
		SELECT COUNT(*) FROM "Event"
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["GET_EVENTS_HOUSE_ALL"], e = Link.Prepare(`
		SELECT e.*, u.login, u.name, s.name, st.short_name, h.name, ht.short_name, n.name, hwt.translate_value
		FROM "Event" AS e
		JOIN "House" AS h ON e.house_id = h.id
        JOIN "Street" AS s ON s.id = h.street_id
        JOIN "Street_type" AS st ON s.type_id = st.id
        JOIN "House_type" AS ht ON h.type_id = ht.id
		JOIN "User" AS u ON e.user_id = u.id
		LEFT JOIN "Node" AS n ON e.node_id = n.id
		LEFT JOIN "Hardware" AS hw ON e.hardware_id = hw.id
		LEFT JOIN "Hardware_type" AS hwt ON hw.type_id = hwt.id
		WHERE e.house_id = $1
		ORDER BY e.created_at DESC
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["GET_EVENTS_HOUSE_ALL_COUNT"], e = Link.Prepare(`
		SELECT COUNT(*) FROM "Event" WHERE house_id = $1
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["GET_EVENTS_HOUSE_ONLY"], e = Link.Prepare(`
		SELECT e.*, u.login, u.name, s.name, st.short_name, h.name, ht.short_name, n.name, hwt.translate_value
		FROM "Event" AS e
		JOIN "House" AS h ON e.house_id = h.id
        JOIN "Street" AS s ON s.id = h.street_id
        JOIN "Street_type" AS st ON s.type_id = st.id
        JOIN "House_type" AS ht ON h.type_id = ht.id
		JOIN "User" AS u ON e.user_id = u.id
		LEFT JOIN "Node" AS n ON e.node_id = n.id
		LEFT JOIN "Hardware" AS hw ON e.hardware_id = hw.id
		LEFT JOIN "Hardware_type" AS hwt ON hw.type_id = hwt.id
		WHERE e.house_id = $1 AND e.node_id IS NULL AND e.hardware_id IS NULL
		ORDER BY e.created_at DESC
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["GET_EVENTS_HOUSE_ONLY_COUNT"], e = Link.Prepare(`
		SELECT COUNT(*) FROM "Event"
		WHERE house_id = $1 AND node_id IS NULL AND hardware_id IS NULL
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["GET_EVENTS_NODE_ALL"], e = Link.Prepare(`
		SELECT e.*, u.login, u.name, s.name, st.short_name, h.name, ht.short_name, n.name, hwt.translate_value
		FROM "Event" AS e
		JOIN "House" AS h ON e.house_id = h.id
        JOIN "Street" AS s ON s.id = h.street_id
        JOIN "Street_type" AS st ON s.type_id = st.id
        JOIN "House_type" AS ht ON h.type_id = ht.id
		JOIN "User" AS u ON e.user_id = u.id
		LEFT JOIN "Node" AS n ON e.node_id = n.id
		LEFT JOIN "Hardware" AS hw ON e.hardware_id = hw.id
		LEFT JOIN "Hardware_type" AS hwt ON hw.type_id = hwt.id
		WHERE e.node_id = $1
		ORDER BY e.created_at DESC
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["GET_EVENTS_NODE_ALL_COUNT"], e = Link.Prepare(`
		SELECT COUNT(*)
		FROM "Event"
		WHERE node_id = $1
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["GET_EVENTS_NODE_ONLY"], e = Link.Prepare(`
		SELECT e.*, u.login, u.name, s.name, st.short_name, h.name, ht.short_name, n.name, hwt.translate_value
		FROM "Event" AS e
		JOIN "House" AS h ON e.house_id = h.id
        JOIN "Street" AS s ON s.id = h.street_id
        JOIN "Street_type" AS st ON s.type_id = st.id
        JOIN "House_type" AS ht ON h.type_id = ht.id
		JOIN "User" AS u ON e.user_id = u.id
		LEFT JOIN "Node" AS n ON e.node_id = n.id
		LEFT JOIN "Hardware" AS hw ON e.hardware_id = hw.id
		LEFT JOIN "Hardware_type" AS hwt ON hw.type_id = hwt.id
		WHERE e.node_id = $1 AND e.hardware_id IS NULL
		ORDER BY e.created_at DESC
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["GET_EVENTS_NODE_ONLY_COUNT"], e = Link.Prepare(`
		SELECT COUNT(*) FROM "Event"
		WHERE node_id = $1 AND hardware_id IS NULL
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["GET_EVENTS_HARDWARE"], e = Link.Prepare(`
		SELECT e.*, u.login, u.name, s.name, st.short_name, h.name, ht.short_name, n.name, hwt.translate_value
		FROM "Event" AS e
		JOIN "House" AS h ON e.house_id = h.id
        JOIN "Street" AS s ON s.id = h.street_id
        JOIN "Street_type" AS st ON s.type_id = st.id
        JOIN "House_type" AS ht ON h.type_id = ht.id
		JOIN "User" AS u ON e.user_id = u.id
		LEFT JOIN "Node" AS n ON e.node_id = n.id
		LEFT JOIN "Hardware" AS hw ON e.hardware_id = hw.id
		LEFT JOIN "Hardware_type" AS hwt ON hw.type_id = hwt.id
		WHERE e.hardware_id = $1
		ORDER BY e.created_at DESC
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["GET_EVENTS_HARDWARE_COUNT"], e = Link.Prepare(`
		SELECT COUNT(*) FROM "Event" WHERE hardware_id = $1
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["CREATE_EVENT"], e = Link.Prepare(`
		INSERT INTO "Event"(house_id, node_id, hardware_id, user_id, description, created_at) VALUES ($1, $2, $3, $4, $5, $6)
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	return errorsList
}

func (event *Event) CreateEvent() error {
	stmt, ok := query["CREATE_EVENT"]
	if !ok {
		err := errors.New("запрос CREATE_EVENT не подготовлен")
		utils.Logger.Println(err)
		return err
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
		event.User.ID,
		event.Description,
		event.CreatedAt,
	)
	if err != nil {
		utils.Logger.Println(err)
		return err
	}

	return nil
}

func GetEvents(from string, id int) ([]Event, int, error) {
	key := "GET_EVENTS"

	if from != "" {
		key += "_" + from
	}

	stmt, ok := query[key]
	if !ok {
		err := errors.New("запрос " + key + " не подготовлен")
		utils.Logger.Println(err)
		return nil, 0, err
	}

	var rows *sql.Rows
	var err error
	var countParam interface{}

	if id > 0 {
		rows, err = stmt.Query(id)
		countParam = id
	} else {
		rows, err = stmt.Query()
	}

	if err != nil {
		utils.Logger.Println(err)
		return nil, 0, err
	}
	defer rows.Close()

	var events []Event

	for rows.Next() {
		var (
			event                      Event
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
			&event.User.ID,
			&event.Description,
			&event.CreatedAt,
			&event.User.Login,
			&event.User.Name,
			&event.Address.Street.Name,
			&event.Address.Street.Type.ShortName,
			&event.Address.House.Name,
			&event.Address.House.Type.ShortName,
			&nodeName,
			&hardwareTypeTranslateValue,
		); err != nil {
			utils.Logger.Println(err)
			return nil, 0, err
		}

		if nodeID.Valid {
			event.Node = &Node{ID: int(nodeID.Int64), Name: nodeName.String}
		}

		if hardwareID.Valid {
			event.Hardware = &Hardware{ID: int(hardwareID.Int64), Type: Reference{TranslateValue: hardwareTypeTranslateValue.String}}
		}

		events = append(events, event)
	}

	count, err := countRecord(key+"_COUNT", countParam)
	if err != nil {
		utils.Logger.Println(err)
		return nil, 0, err
	}

	return events, count, nil
}
