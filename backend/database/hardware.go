package database

import (
	"backend/utils"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
)

type Switch struct {
	ID               int
	Name             string
	OperationMode    Enum
	PortAmount       int
	CommunityRead    sql.NullString
	CommunityWrite   sql.NullString
	FirmwareOID      sql.NullString
	SystemNameOID    sql.NullString
	SerialNumberOID  sql.NullString
	SaveConfigOID    sql.NullString
	PortDescOID      sql.NullString
	VlanOID          sql.NullString
	PortUntaggedOID  sql.NullString
	SpeedOID         sql.NullString
	BatteryStatusOID sql.NullString
	BatteryChargeOID sql.NullString
	PortModeOID      sql.NullString
	UptimeOID        sql.NullString
	CreatedAt        int64
}

type Hardware struct {
	ID          int
	Node        Node
	Type        Enum
	Switch      Switch
	IpAddress   sql.NullString
	MgmtVlan    sql.NullString
	Description sql.NullString
	CreatedAt   int64
	UpdatedAt   sql.NullInt64
}

func prepareHardware() []string {
	var e error
	errorsList := make([]string, 0)

	if query == nil {
		query = make(map[string]*sql.Stmt)
	}

	query["GET_HARDWARE_TYPES"], e = Link.Prepare(`
		SELECT * FROM "Hardware_type" ORDER BY id
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["GET_OPERATION_MODES"], e = Link.Prepare(`
		SELECT * FROM "Operation_mode" ORDER BY id
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["CREATE_SWITCH"], e = Link.Prepare(`
		INSERT INTO "Switch"(name, operation_mode_id, community_read, community_write, port_amount, firmware_oid, 
		                     system_name_oid, sn_oid, save_config_oid, port_desc_oid, vlan_oid, port_untagged_oid, 
		                     speed_oid, battery_status_oid, battery_charge_oid, port_mode_oid, uptime_oid, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
		RETURNING id
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["EDIT_SWITCH"], e = Link.Prepare(`
		UPDATE "Switch" SET name = $2, operation_mode_id = $3, community_read = $4, community_write = $5, port_amount = $6,
		                    firmware_oid = $7, system_name_oid = $8, sn_oid = $9, save_config_oid = $10, port_desc_oid = $11,
		                    vlan_oid = $12, port_untagged_oid = $13, speed_oid = $14, battery_status_oid = $15, battery_charge_oid = $16,
		                    port_mode_oid = $17, uptime_oid = $18
		WHERE id = $1
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["GET_SWITCHES"], e = Link.Prepare(`
		SELECT s.*, om.value, om.translate_value 
		FROM "Switch" AS s
		LEFT JOIN "Operation_mode" AS om ON s.operation_mode_id = om.id
		ORDER BY s.id
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["GET_HARDWARE"], e = Link.Prepare(`
		SELECT hd.*, s.name, st.short_name, h.name, ht.short_name, hdt.value, hdt.translate_value, sw.name, n.name
		FROM "Hardware" AS hd
		JOIN "Node" AS n ON hd.node_id = n.id
		JOIN "House" AS h ON n.house_id = h.id
        JOIN "Street" AS s ON s.id = h.street_id
        JOIN "Street_type" AS st ON s.type_id = st.id
        JOIN "House_type" AS ht ON h.type_id = ht.id
		JOIN "Hardware_type" AS hdt ON hd.type_id = hdt.id
		LEFT JOIN "Switch" AS sw ON hd.switch_id = sw.id
		ORDER BY hd.id
		OFFSET $1
		LIMIT 20
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["GET_HARDWARE_COUNT"], e = Link.Prepare(`
		SELECT COUNT(*) FROM "Hardware"
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["GET_HOUSE_HARDWARE"], e = Link.Prepare(`
		SELECT hd.*, s.name, st.short_name, h.name, ht.short_name, hdt.value, hdt.translate_value, sw.name, n.name
		FROM "Hardware" AS hd
		JOIN "Node" AS n ON hd.node_id = n.id
		JOIN "House" AS h ON n.house_id = h.id
        JOIN "Street" AS s ON s.id = h.street_id
        JOIN "Street_type" AS st ON s.type_id = st.id
        JOIN "House_type" AS ht ON h.type_id = ht.id
		JOIN "Hardware_type" AS hdt ON hd.type_id = hdt.id
		LEFT JOIN "Switch" AS sw ON hd.switch_id = sw.id
		WHERE n.house_id = $1
		ORDER BY hd.id
		OFFSET $2
		LIMIT 20
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["GET_HOUSE_HARDWARE_COUNT"], e = Link.Prepare(`
		SELECT COUNT(hd.*) 
		FROM "Hardware" AS hd 
		JOIN "Node" AS n ON hd.node_id = n.id
		WHERE n.house_id = $1
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["GET_NODE_HARDWARE"], e = Link.Prepare(`
		SELECT hd.*, s.name, st.short_name, h.name, ht.short_name, hdt.value, hdt.translate_value, sw.name, n.name
		FROM "Hardware" AS hd
		JOIN "Node" AS n ON hd.node_id = n.id
		JOIN "House" AS h ON n.house_id = h.id
        JOIN "Street" AS s ON s.id = h.street_id
        JOIN "Street_type" AS st ON s.type_id = st.id
        JOIN "House_type" AS ht ON h.type_id = ht.id
		JOIN "Hardware_type" AS hdt ON hd.type_id = hdt.id
		LEFT JOIN "Switch" AS sw ON hd.switch_id = sw.id
		WHERE hd.node_id = $1
		ORDER BY hd.id
		OFFSET $2
		LIMIT 20
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["GET_NODE_HARDWARE_COUNT"], e = Link.Prepare(`
		SELECT COUNT(*) 
		FROM "Hardware"
		WHERE node_id = $1
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["GET_SEARCH_HARDWARE"], e = Link.Prepare(`
		SELECT hd.*, s.name, st.short_name, h.name, ht.short_name, hdt.value, hdt.translate_value, sw.name, n.name
		FROM "Hardware" AS hd
		JOIN "Node" AS n ON hd.node_id = n.id
		JOIN "House" AS h ON n.house_id = h.id
        JOIN "Street" AS s ON s.id = h.street_id
        JOIN "Street_type" AS st ON s.type_id = st.id
        JOIN "House_type" AS ht ON h.type_id = ht.id
		JOIN "Hardware_type" AS hdt ON hd.type_id = hdt.id
		LEFT JOIN "Switch" AS sw ON hd.switch_id = sw.id
		WHERE n.name ILIKE '%' || $1 || '%'
			OR hdt.translate_value ILIKE '%' || $1 || '%'
			OR sw.name ILIKE '%' || $1 || '%'
			OR hd.ip_address ILIKE '%' || $1 || '%'
			OR s.name ILIKE '%' || $1 || '%'
			OR st.short_name ILIKE '%' || $1 || '%'
			OR h.name ILIKE '%' || $1 || '%'
			OR ht.short_name ILIKE '%' || $1 || '%'
			OR (st.short_name || ' ' || s.name || ', ' || ht.short_name || ' ' || h.name) ILIKE '%' || $1 || '%'
		ORDER BY hd.id
		OFFSET $2
		LIMIT 20
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["GET_SEARCH_HARDWARE_COUNT"], e = Link.Prepare(`
		SELECT COUNT(hd.*)
		FROM "Hardware" AS hd
		JOIN "Node" AS n ON hd.node_id = n.id
		JOIN "House" AS h ON n.house_id = h.id
        JOIN "Street" AS s ON s.id = h.street_id
        JOIN "Street_type" AS st ON s.type_id = st.id
        JOIN "House_type" AS ht ON h.type_id = ht.id
		JOIN "Hardware_type" AS hdt ON hd.type_id = hdt.id
		LEFT JOIN "Switch" AS sw ON hd.switch_id = sw.id
		WHERE n.name ILIKE '%' || $1 || '%'
			OR hdt.translate_value ILIKE '%' || $1 || '%'
			OR sw.name ILIKE '%' || $1 || '%'
			OR hd.ip_address ILIKE '%' || $1 || '%'
			OR s.name ILIKE '%' || $1 || '%'
			OR st.short_name ILIKE '%' || $1 || '%'
			OR h.name ILIKE '%' || $1 || '%'
			OR ht.short_name ILIKE '%' || $1 || '%'
			OR (st.short_name || ' ' || s.name || ', ' || ht.short_name || ' ' || h.name) ILIKE '%' || $1 || '%'
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["CREATE_HARDWARE"], e = Link.Prepare(`
		INSERT INTO "Hardware" (node_id, type_id, switch_id, ip_address, mgmt_vlan, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["EDIT_HARDWARE"], e = Link.Prepare(`
		UPDATE "Hardware" SET node_id = $2, type_id = $3, switch_id = $4, ip_address = $5, mgmt_vlan = $6, description = $7,
		                      updated_at = $8
		WHERE id = $1
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["GET_HARDWARE_BY_ID"], e = Link.Prepare(`
		SELECT hd.*, s.name, st.short_name, h.name, ht.short_name, hdt.value, hdt.translate_value, sw.name, n.name
		FROM "Hardware" AS hd
		JOIN "Node" AS n ON hd.node_id = n.id
		JOIN "House" AS h ON n.house_id = h.id
        JOIN "Street" AS s ON s.id = h.street_id
        JOIN "Street_type" AS st ON s.type_id = st.id
        JOIN "House_type" AS ht ON h.type_id = ht.id
		JOIN "Hardware_type" AS hdt ON hd.type_id = hdt.id
		LEFT JOIN "Switch" AS sw ON hd.switch_id = sw.id
		WHERE hd.id = $1
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["GET_HARDWARE_FILES"], e = Link.Prepare(`
		SELECT * FROM "Hardware_files" WHERE hardware_id = $1
		ORDER BY upload_at DESC
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	return errorsList
}

func GetHardwareFiles(hardwareID int) ([]File, error) {
	stmt, ok := query["GET_HARDWARE_FILES"]
	if !ok {
		err := errors.New("запрос GET_HARDWARE_FILES не подготовлен")
		utils.Logger.Println(err)
		return nil, err
	}

	rows, err := stmt.Query(hardwareID)
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
			&file.Hardware.ID,
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

func (hardware *Hardware) GetHardwareByID() error {
	stmt, ok := query["GET_HARDWARE_BY_ID"]
	if !ok {
		err := errors.New("запрос GET_HARDWARE_BY_ID не подготовлен")
		utils.Logger.Println(err)
		return err
	}

	var (
		switchID   sql.NullInt64
		switchName sql.NullString
	)

	if err := stmt.QueryRow(hardware.ID).Scan(
		&hardware.ID,
		&hardware.Node.ID,
		&hardware.Type.ID,
		&switchID,
		&hardware.IpAddress,
		&hardware.MgmtVlan,
		&hardware.Description,
		&hardware.CreatedAt,
		&hardware.UpdatedAt,
		&hardware.Node.Address.Street.Name,
		&hardware.Node.Address.Street.Type.ShortName,
		&hardware.Node.Address.House.Name,
		&hardware.Node.Address.House.Type.ShortName,
		&hardware.Type.Value,
		&hardware.Type.TranslateValue,
		&switchName,
		&hardware.Node.Name,
	); err != nil {
		utils.Logger.Println(err)
		return err
	}

	if switchID.Valid {
		hardware.Switch = Switch{ID: int(switchID.Int64), Name: switchName.String}
	}

	return nil
}

func (hardware *Hardware) EditHardware() error {
	stmt, ok := query["EDIT_HARDWARE"]
	if !ok {
		err := errors.New("запрос EDIT_HARDWARE не подготовлен")
		utils.Logger.Println(err)
		return err
	}

	var switchID interface{}

	if hardware.Switch.ID != 0 {
		switchID = hardware.Switch.ID
	}

	_, err := stmt.Exec(
		hardware.ID,
		hardware.Node.ID,
		hardware.Type.ID,
		switchID,
		hardware.IpAddress,
		hardware.MgmtVlan,
		hardware.Description,
		hardware.UpdatedAt,
	)
	if err != nil {
		utils.Logger.Println(err)
		return err
	}

	return nil
}

func (hardware *Hardware) CreateHardware() error {
	stmt, ok := query["CREATE_HARDWARE"]
	if !ok {
		err := errors.New("запрос CREATE_HARDWARE не подготовлен")
		utils.Logger.Println(err)
		return err
	}

	var switchID interface{}

	if hardware.Switch.ID != 0 {
		switchID = hardware.Switch.ID
	}

	if err := stmt.QueryRow(
		hardware.Node.ID,
		hardware.Type.ID,
		switchID,
		hardware.IpAddress,
		hardware.MgmtVlan,
		hardware.Description,
		hardware.CreatedAt,
		nil,
	).Scan(&hardware.ID); err != nil {
		utils.Logger.Println(err)
		return err
	}

	return nil
}

func GetSearchHardware(search string, offset int) ([]Hardware, int, error) {
	stmt, ok := query["GET_SEARCH_HARDWARE"]
	if !ok {
		err := errors.New("запрос GET_SEARCH_HARDWARE не подготовлен")
		utils.Logger.Println(err)
		return nil, 0, err
	}

	count, err := countRecord("GET_SEARCH_HARDWARE_COUNT", search)
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

	var hardware []Hardware

	for rows.Next() {
		var (
			_hardware  Hardware
			switchID   sql.NullInt64
			switchName sql.NullString
		)

		if err = rows.Scan(
			&_hardware.ID,
			&_hardware.Node.ID,
			&_hardware.Type.ID,
			&switchID,
			&_hardware.IpAddress,
			&_hardware.MgmtVlan,
			&_hardware.Description,
			&_hardware.CreatedAt,
			&_hardware.UpdatedAt,
			&_hardware.Node.Address.Street.Name,
			&_hardware.Node.Address.Street.Type.ShortName,
			&_hardware.Node.Address.House.Name,
			&_hardware.Node.Address.House.Type.ShortName,
			&_hardware.Type.Value,
			&_hardware.Type.TranslateValue,
			&switchName,
			&_hardware.Node.Name,
		); err != nil && !errors.Is(err, sql.ErrNoRows) {
			utils.Logger.Println(err)
			return nil, 0, err
		}

		if switchID.Valid {
			_hardware.Switch = Switch{ID: int(switchID.Int64), Name: switchName.String}
		}

		hardware = append(hardware, _hardware)
	}

	return hardware, count, nil
}

func GetNodeHardware(nodeID int, offset int) ([]Hardware, int, error) {
	stmt, ok := query["GET_NODE_HARDWARE"]
	if !ok {
		err := errors.New("запрос GET_NODE_HARDWARE не подготовлен")
		utils.Logger.Println(err)
		return nil, 0, err
	}

	count, err := countRecord("GET_HOUSE_HARDWARE_COUNT", nodeID)
	if err != nil {
		utils.Logger.Println(err)
		return nil, 0, err
	}

	rows, err := stmt.Query(nodeID, offset)
	if err != nil {
		utils.Logger.Println(err)
		return nil, 0, err
	}
	defer rows.Close()

	var hardware []Hardware

	for rows.Next() {
		var (
			_hardware  Hardware
			switchID   sql.NullInt64
			switchName sql.NullString
		)

		if err = rows.Scan(
			&_hardware.ID,
			&_hardware.Node.ID,
			&_hardware.Type.ID,
			&switchID,
			&_hardware.IpAddress,
			&_hardware.MgmtVlan,
			&_hardware.Description,
			&_hardware.CreatedAt,
			&_hardware.UpdatedAt,
			&_hardware.Node.Address.Street.Name,
			&_hardware.Node.Address.Street.Type.ShortName,
			&_hardware.Node.Address.House.Name,
			&_hardware.Node.Address.House.Type.ShortName,
			&_hardware.Type.Value,
			&_hardware.Type.TranslateValue,
			&switchName,
			&_hardware.Node.Name,
		); err != nil && !errors.Is(err, sql.ErrNoRows) {
			utils.Logger.Println(err)
			return nil, 0, err
		}

		if switchID.Valid {
			_hardware.Switch = Switch{ID: int(switchID.Int64), Name: switchName.String}
		}

		hardware = append(hardware, _hardware)
	}

	return hardware, count, nil
}

func GetHouseHardware(houseID int, offset int) ([]Hardware, int, error) {
	stmt, ok := query["GET_HOUSE_HARDWARE"]
	if !ok {
		err := errors.New("запрос GET_HOUSE_HARDWARE не подготовлен")
		utils.Logger.Println(err)
		return nil, 0, err
	}

	count, err := countRecord("GET_HOUSE_HARDWARE_COUNT", houseID)
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

	var hardware []Hardware

	for rows.Next() {
		var (
			_hardware  Hardware
			switchID   sql.NullInt64
			switchName sql.NullString
		)

		if err = rows.Scan(
			&_hardware.ID,
			&_hardware.Node.ID,
			&_hardware.Type.ID,
			&switchID,
			&_hardware.IpAddress,
			&_hardware.MgmtVlan,
			&_hardware.Description,
			&_hardware.CreatedAt,
			&_hardware.UpdatedAt,
			&_hardware.Node.Address.Street.Name,
			&_hardware.Node.Address.Street.Type.ShortName,
			&_hardware.Node.Address.House.Name,
			&_hardware.Node.Address.House.Type.ShortName,
			&_hardware.Type.Value,
			&_hardware.Type.TranslateValue,
			&switchName,
			&_hardware.Node.Name,
		); err != nil && !errors.Is(err, sql.ErrNoRows) {
			utils.Logger.Println(err)
			return nil, 0, err
		}

		if switchID.Valid {
			_hardware.Switch = Switch{ID: int(switchID.Int64), Name: switchName.String}
		}

		hardware = append(hardware, _hardware)
	}

	return hardware, count, nil
}

func GetHardware(offset int) ([]Hardware, int, error) {
	stmt, ok := query["GET_HARDWARE"]
	if !ok {
		err := errors.New("запрос GET_HARDWARE не подготовлен")
		utils.Logger.Println(err)
		return nil, 0, err
	}

	count, err := countRecord("GET_HARDWARE_COUNT", nil)
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

	var hardware []Hardware

	for rows.Next() {
		var (
			_hardware  Hardware
			switchID   sql.NullInt64
			switchName sql.NullString
		)

		if err = rows.Scan(
			&_hardware.ID,
			&_hardware.Node.ID,
			&_hardware.Type.ID,
			&switchID,
			&_hardware.IpAddress,
			&_hardware.MgmtVlan,
			&_hardware.Description,
			&_hardware.CreatedAt,
			&_hardware.UpdatedAt,
			&_hardware.Node.Address.Street.Name,
			&_hardware.Node.Address.Street.Type.ShortName,
			&_hardware.Node.Address.House.Name,
			&_hardware.Node.Address.House.Type.ShortName,
			&_hardware.Type.Value,
			&_hardware.Type.TranslateValue,
			&switchName,
			&_hardware.Node.Name,
		); err != nil && !errors.Is(err, sql.ErrNoRows) {
			utils.Logger.Println(err)
			return nil, 0, err
		}

		if switchID.Valid {
			_hardware.Switch = Switch{ID: int(switchID.Int64), Name: switchName.String}
		}

		hardware = append(hardware, _hardware)
	}

	return hardware, count, nil
}

func GetSwitches() ([]Switch, error) {
	stmt, ok := query["GET_SWITCHES"]
	if !ok {
		err := errors.New("запрос GET_SWITCHES не подготовлен")
		utils.Logger.Println(err)
		return nil, err
	}

	rows, err := stmt.Query()
	if err != nil {
		utils.Logger.Println(err)
		return nil, err
	}
	defer rows.Close()

	var switches []Switch

	for rows.Next() {
		var (
			operationModeID             sql.NullInt64
			operationModeValue          sql.NullString
			operationModeTranslateValue sql.NullString
			_switch                     Switch
		)
		if err = rows.Scan(
			&_switch.ID,
			&_switch.Name,
			&operationModeID,
			&_switch.CommunityRead,
			&_switch.CommunityWrite,
			&_switch.PortAmount,
			&_switch.FirmwareOID,
			&_switch.SystemNameOID,
			&_switch.SerialNumberOID,
			&_switch.SaveConfigOID,
			&_switch.PortDescOID,
			&_switch.VlanOID,
			&_switch.PortUntaggedOID,
			&_switch.SpeedOID,
			&_switch.BatteryStatusOID,
			&_switch.BatteryChargeOID,
			&_switch.PortModeOID,
			&_switch.UptimeOID,
			&_switch.CreatedAt,
			&operationModeValue,
			&operationModeTranslateValue,
		); err != nil {
			utils.Logger.Println(err)
			return nil, err
		}

		if operationModeID.Valid {
			_switch.OperationMode = Enum{
				ID:             int(operationModeID.Int64),
				Value:          operationModeValue.String,
				TranslateValue: operationModeTranslateValue.String,
			}
		}

		switches = append(switches, _switch)
	}

	return switches, nil
}

func (_switch *Switch) EditSwitch() error {
	stmt, ok := query["EDIT_SWITCH"]
	if !ok {
		err := errors.New("запрос EDIT_SWITCH не подготовлен")
		utils.Logger.Println(err)
		return err
	}

	var operationModeID interface{}

	if _switch.OperationMode.ID != 0 {
		operationModeID = _switch.OperationMode.ID
	}

	_, err := stmt.Exec(
		_switch.ID,
		_switch.Name,
		operationModeID,
		_switch.CommunityRead,
		_switch.CommunityWrite,
		_switch.PortAmount,
		_switch.FirmwareOID,
		_switch.SystemNameOID,
		_switch.SerialNumberOID,
		_switch.SaveConfigOID,
		_switch.PortDescOID,
		_switch.VlanOID,
		_switch.PortUntaggedOID,
		_switch.SpeedOID,
		_switch.BatteryStatusOID,
		_switch.BatteryChargeOID,
		_switch.PortModeOID,
		_switch.UptimeOID,
	)
	if err != nil {
		utils.Logger.Println(err)
		return err
	}

	return nil
}

func (_switch *Switch) CreateSwitch() error {
	stmt, ok := query["CREATE_SWITCH"]
	if !ok {
		err := errors.New("запрос CREATE_SWITCH не подготовлен")
		utils.Logger.Println(err)
		return err
	}

	var operationModeID interface{}

	if _switch.OperationMode.ID != 0 {
		operationModeID = _switch.OperationMode.ID
	}

	if err := stmt.QueryRow(
		_switch.Name,
		operationModeID,
		_switch.CommunityRead,
		_switch.CommunityWrite,
		_switch.PortAmount,
		_switch.FirmwareOID,
		_switch.SystemNameOID,
		_switch.SerialNumberOID,
		_switch.SaveConfigOID,
		_switch.PortDescOID,
		_switch.VlanOID,
		_switch.PortUntaggedOID,
		_switch.SpeedOID,
		_switch.BatteryStatusOID,
		_switch.BatteryChargeOID,
		_switch.PortModeOID,
		_switch.UptimeOID,
		_switch.CreatedAt,
	).Scan(&_switch.ID); err != nil {
		utils.Logger.Println(err)
		return err
	}

	return nil
}

func GetHardwareTypes() ([]Enum, error) {
	stmt, ok := query["GET_HARDWARE_TYPES"]
	if !ok {
		err := errors.New("запрос GET_HARDWARE_TYPES не подготовлен")
		utils.Logger.Println(err)
		return nil, err
	}

	rows, err := stmt.Query()
	if err != nil {
		utils.Logger.Println(err)
		return nil, err
	}
	defer rows.Close()

	var hardwareTypes []Enum
	for rows.Next() {
		var hardwareType Enum
		if err = rows.Scan(
			&hardwareType.ID,
			&hardwareType.Value,
			&hardwareType.TranslateValue,
			&hardwareType.CreatedAt,
		); err != nil {
			utils.Logger.Println(err)
			return nil, err
		}

		hardwareTypes = append(hardwareTypes, hardwareType)
	}

	return hardwareTypes, nil
}

func GetOperationModes() ([]Enum, error) {
	stmt, ok := query["GET_OPERATION_MODES"]
	if !ok {
		err := errors.New("запрос GET_OPERATION_MODES не подготовлен")
		utils.Logger.Println(err)
		return nil, err
	}

	rows, err := stmt.Query()
	if err != nil {
		utils.Logger.Println(err)
		return nil, err
	}
	defer rows.Close()

	var operationModes []Enum
	for rows.Next() {
		var operationMode Enum
		if err = rows.Scan(
			&operationMode.ID,
			&operationMode.Value,
			&operationMode.TranslateValue,
			&operationMode.CreatedAt,
		); err != nil {
			utils.Logger.Println(err)
			return nil, err
		}

		operationModes = append(operationModes, operationMode)
	}

	return operationModes, nil
}

func (hardware *Hardware) ValidateHardware() bool {
	fmt.Println(hardware)

	if hardware.Type.ID == 0 || hardware.Node.ID == 0 {
		return false
	}

	if hardware.Type.Value == "switch" && (hardware.Switch.ID == 0 || !hardware.IpAddress.Valid) {
		return false
	}

	return true
}
