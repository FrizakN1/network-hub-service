package database

import (
	"backend/models"
	"database/sql"
	"errors"
)

type SwitchRepository interface {
	GetSwitches() ([]models.Switch, error)
	EditSwitch(_switch *models.Switch) error
	CreateSwitch(_switch *models.Switch) error
}

type DefaultSwitchRepository struct {
	Database Database
}

func (r *DefaultSwitchRepository) GetSwitches() ([]models.Switch, error) {
	stmt, ok := r.Database.GetQuery("GET_SWITCHES")
	if !ok {
		return nil, errors.New("query GET_SWITCHES is not prepare")
	}

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var switches []models.Switch

	for rows.Next() {
		var (
			operationModeID    sql.NullInt64
			operationModeKey   sql.NullString
			operationModeValue sql.NullString
			_switch            models.Switch
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
			&operationModeKey,
			&operationModeValue,
		); err != nil {
			return nil, err
		}

		if operationModeID.Valid {
			_switch.OperationMode = models.Reference{
				ID:    int(operationModeID.Int64),
				Key:   operationModeKey.String,
				Value: operationModeValue.String,
			}
		}

		switches = append(switches, _switch)
	}

	return switches, nil
}

func (r *DefaultSwitchRepository) EditSwitch(_switch *models.Switch) error {
	stmt, ok := r.Database.GetQuery("EDIT_SWITCH")
	if !ok {
		return errors.New("query EDIT_SWITCH is not prepare")
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
		return err
	}

	return nil
}

func (r *DefaultSwitchRepository) CreateSwitch(_switch *models.Switch) error {
	stmt, ok := r.Database.GetQuery("CREATE_SWITCH")
	if !ok {
		return errors.New("query CREATE_SWITCH is not prepare")
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
		return err
	}

	return nil
}
