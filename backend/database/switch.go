package database

import (
	"backend/utils"
	"database/sql"
	"errors"
)

type Switch struct {
	ID               int
	Name             string
	OperationMode    Reference
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
	MacOID           sql.NullString
}

func prepareSwitch() []string {
	var e error
	errorsList := make([]string, 0)

	if query == nil {
		query = make(map[string]*sql.Stmt)
	}

	query["CREATE_SWITCH"], e = Link.Prepare(`
		INSERT INTO "Switch"(name, operation_mode_id, community_read, community_write, port_amount, firmware_oid, 
		                     system_name_oid, sn_oid, save_config_oid, port_desc_oid, vlan_oid, port_untagged_oid, 
		                     speed_oid, battery_status_oid, battery_charge_oid, port_mode_oid, uptime_oid, created_at, mac_oid) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
		RETURNING id
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["EDIT_SWITCH"], e = Link.Prepare(`
		UPDATE "Switch" SET name = $2, operation_mode_id = $3, community_read = $4, community_write = $5, port_amount = $6,
		                    firmware_oid = $7, system_name_oid = $8, sn_oid = $9, save_config_oid = $10, port_desc_oid = $11,
		                    vlan_oid = $12, port_untagged_oid = $13, speed_oid = $14, battery_status_oid = $15, battery_charge_oid = $16,
		                    port_mode_oid = $17, uptime_oid = $18, mac_oid = $19
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

	return errorsList
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
			&_switch.MacOID,
			&operationModeValue,
			&operationModeTranslateValue,
		); err != nil {
			utils.Logger.Println(err)
			return nil, err
		}

		if operationModeID.Valid {
			_switch.OperationMode = Reference{
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
		_switch.MacOID,
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
		_switch.MacOID,
	).Scan(&_switch.ID); err != nil {
		utils.Logger.Println(err)
		return err
	}

	return nil
}
