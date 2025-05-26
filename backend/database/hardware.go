package database

import (
	"backend/models"
	"database/sql"
	"errors"
)

type HardwareRepository interface {
	GetHardwareByID(hardware *models.Hardware) error
	EditHardware(hardware *models.Hardware) error
	CreateHardware(hardware *models.Hardware) error
	GetSearchHardware(search string, offset int) ([]models.Hardware, int, error)
	GetHardware(offset int, houseID int, nodeID int) ([]models.Hardware, int, error)
	ValidateHardware(hardware models.Hardware) bool
	DeleteHardware(hardwareID int) error
}

type DefaultHardwareRepository struct {
	Database Database
}

func (r *DefaultHardwareRepository) DeleteHardware(hardwareID int) error {
	stmt, ok := r.Database.GetQuery("DELETE_HARDWARE")
	if !ok {
		return errors.New("query DELETE_HARDWARE is not prepare")
	}

	_, err := stmt.Exec(hardwareID)
	if err != nil {
		return err
	}

	return nil
}

func (r *DefaultHardwareRepository) GetHardwareByID(hardware *models.Hardware) error {
	stmt, ok := r.Database.GetQuery("GET_HARDWARE_BY_ID")
	if !ok {
		return errors.New("query GET_HARDWARE_BY_ID is not prepare")
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
		&hardware.IsDelete,
		&hardware.Node.Address.Street.Name,
		&hardware.Node.Address.Street.Type.ShortName,
		&hardware.Node.Address.House.ID,
		&hardware.Node.Address.House.Name,
		&hardware.Node.Address.House.Type.ShortName,
		&hardware.Type.Key,
		&hardware.Type.Value,
		&switchName,
		&hardware.Node.Name,
	); err != nil {
		return err
	}

	if switchID.Valid {
		hardware.Switch = models.Switch{ID: int(switchID.Int64), Name: switchName.String}
	}

	return nil
}

func (r *DefaultHardwareRepository) EditHardware(hardware *models.Hardware) error {
	stmt, ok := r.Database.GetQuery("EDIT_HARDWARE")
	if !ok {
		return errors.New("query EDIT_HARDWARE is not prepare")
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
		return err
	}

	return nil
}

func (r *DefaultHardwareRepository) CreateHardware(hardware *models.Hardware) error {
	stmt, ok := r.Database.GetQuery("CREATE_HARDWARE")
	if !ok {
		return errors.New("query CREATE_HARDWARE is not prepare")
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
		return err
	}

	return nil
}

func (r *DefaultHardwareRepository) GetSearchHardware(search string, offset int) ([]models.Hardware, int, error) {
	stmt, ok := r.Database.GetQuery("GET_SEARCH_HARDWARE")
	if !ok {
		return nil, 0, errors.New("query GET_SEARCH_HARDWARE is not prepare")
	}

	rows, err := stmt.Query(search, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var hardware []models.Hardware
	var count int

	for rows.Next() {
		var (
			_hardware  models.Hardware
			switchID   sql.NullInt64
			switchName sql.NullString
		)

		if err = rows.Scan(
			&_hardware.ID,
			&_hardware.Node.ID,
			&_hardware.Type.ID,
			&switchID,
			&_hardware.IpAddress,
			&_hardware.Node.Address.Street.Name,
			&_hardware.Node.Address.Street.Type.ShortName,
			&_hardware.Node.Address.House.ID,
			&_hardware.Node.Address.House.Name,
			&_hardware.Node.Address.House.Type.ShortName,
			&_hardware.Type.Key,
			&_hardware.Type.Value,
			&switchName,
			&_hardware.Node.Name,
			&count,
		); err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, 0, err
		}

		if switchID.Valid {
			_hardware.Switch = models.Switch{ID: int(switchID.Int64), Name: switchName.String}
		}

		hardware = append(hardware, _hardware)
	}

	return hardware, count, nil
}

func (r *DefaultHardwareRepository) GetHardware(offset int, houseID int, nodeID int) ([]models.Hardware, int, error) {
	stmt, ok := r.Database.GetQuery("GET_HARDWARE")
	if !ok {
		return nil, 0, errors.New("query GET_HARDWARE is not prepare")
	}

	rows, err := stmt.Query(offset, houseID, nodeID)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var hardware []models.Hardware
	var count int

	for rows.Next() {
		var (
			_hardware  models.Hardware
			switchID   sql.NullInt64
			switchName sql.NullString
		)

		if err = rows.Scan(
			&_hardware.ID,
			&_hardware.Node.ID,
			&_hardware.Type.ID,
			&switchID,
			&_hardware.IpAddress,
			&_hardware.Node.Address.Street.Name,
			&_hardware.Node.Address.Street.Type.ShortName,
			&_hardware.Node.Address.House.ID,
			&_hardware.Node.Address.House.Name,
			&_hardware.Node.Address.House.Type.ShortName,
			&_hardware.Type.Key,
			&_hardware.Type.Value,
			&switchName,
			&_hardware.Node.Name,
			&count,
		); err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, 0, err
		}

		if switchID.Valid {
			_hardware.Switch = models.Switch{ID: int(switchID.Int64), Name: switchName.String}
		}

		hardware = append(hardware, _hardware)
	}

	return hardware, count, nil
}

func (r *DefaultHardwareRepository) ValidateHardware(hardware models.Hardware) bool {
	if hardware.Type.ID == 0 || hardware.Node.ID == 0 || hardware.Node.IsPassive {
		return false
	}

	if hardware.Type.Value == "switch" && (hardware.Switch.ID == 0 || !hardware.IpAddress.Valid) {
		return false
	}

	return true
}
