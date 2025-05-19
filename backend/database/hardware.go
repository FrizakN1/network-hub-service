package database

import (
	"backend/models"
	"database/sql"
	"errors"
	"fmt"
)

type HardwareRepository interface {
	GetHardwareByID(hardware *models.Hardware) error
	EditHardware(hardware *models.Hardware) error
	CreateHardware(hardware *models.Hardware) error
	GetSearchHardware(search string, offset int) ([]models.Hardware, int, error)
	GetNodeHardware(nodeID int, offset int) ([]models.Hardware, int, error)
	GetHouseHardware(houseID int, offset int) ([]models.Hardware, int, error)
	GetHardware(offset int) ([]models.Hardware, int, error)
	ValidateHardware(hardware models.Hardware) bool
	DeleteHardware(hardwareID int) error
}

type DefaultHardwareRepository struct {
	Database Database
	Counter  Counter
}

func (r *DefaultHardwareRepository) DeleteHardware(hardwareID int) error {
	stmt, ok := r.Database.GetQuery("DELETE_HARDWARE")
	if !ok {
		return errors.New("запрос DELETE_HARDWARE не подготовлен")
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
		return errors.New("запрос GET_HARDWARE_BY_ID не подготовлен")
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
		&hardware.Type.Value,
		&hardware.Type.TranslateValue,
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
		return errors.New("запрос EDIT_HARDWARE не подготовлен")
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
		return errors.New("запрос CREATE_HARDWARE не подготовлен")
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
		return nil, 0, errors.New("запрос GET_SEARCH_HARDWARE не подготовлен")
	}

	count, err := r.Counter.countRecords("GET_SEARCH_HARDWARE_COUNT", []interface{}{search})
	if err != nil {
		return nil, 0, err
	}

	rows, err := stmt.Query(search, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var hardware []models.Hardware

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
			&_hardware.MgmtVlan,
			&_hardware.Description,
			&_hardware.CreatedAt,
			&_hardware.UpdatedAt,
			&_hardware.IsDelete,
			&_hardware.Node.Address.Street.Name,
			&_hardware.Node.Address.Street.Type.ShortName,
			&_hardware.Node.Address.House.Name,
			&_hardware.Node.Address.House.Type.ShortName,
			&_hardware.Type.Value,
			&_hardware.Type.TranslateValue,
			&switchName,
			&_hardware.Node.Name,
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

func (r *DefaultHardwareRepository) GetNodeHardware(nodeID int, offset int) ([]models.Hardware, int, error) {
	stmt, ok := r.Database.GetQuery("GET_NODE_HARDWARE")
	if !ok {
		return nil, 0, errors.New("запрос GET_NODE_HARDWARE не подготовлен")
	}

	count, err := r.Counter.countRecords("GET_HOUSE_HARDWARE_COUNT", []interface{}{nodeID})
	if err != nil {
		return nil, 0, err
	}

	rows, err := stmt.Query(nodeID, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var hardware []models.Hardware

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
			&_hardware.MgmtVlan,
			&_hardware.Description,
			&_hardware.CreatedAt,
			&_hardware.UpdatedAt,
			&_hardware.IsDelete,
			&_hardware.Node.Address.Street.Name,
			&_hardware.Node.Address.Street.Type.ShortName,
			&_hardware.Node.Address.House.Name,
			&_hardware.Node.Address.House.Type.ShortName,
			&_hardware.Type.Value,
			&_hardware.Type.TranslateValue,
			&switchName,
			&_hardware.Node.Name,
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

func (r *DefaultHardwareRepository) GetHouseHardware(houseID int, offset int) ([]models.Hardware, int, error) {
	stmt, ok := r.Database.GetQuery("GET_HOUSE_HARDWARE")
	if !ok {
		return nil, 0, errors.New("запрос GET_HOUSE_HARDWARE не подготовлен")
	}

	count, err := r.Counter.countRecords("GET_HOUSE_HARDWARE_COUNT", []interface{}{houseID})
	if err != nil {
		return nil, 0, err
	}

	rows, err := stmt.Query(houseID, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var hardware []models.Hardware

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
			&_hardware.MgmtVlan,
			&_hardware.Description,
			&_hardware.CreatedAt,
			&_hardware.UpdatedAt,
			&_hardware.IsDelete,
			&_hardware.Node.Address.Street.Name,
			&_hardware.Node.Address.Street.Type.ShortName,
			&_hardware.Node.Address.House.Name,
			&_hardware.Node.Address.House.Type.ShortName,
			&_hardware.Type.Value,
			&_hardware.Type.TranslateValue,
			&switchName,
			&_hardware.Node.Name,
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

func (r *DefaultHardwareRepository) GetHardware(offset int) ([]models.Hardware, int, error) {
	stmt, ok := r.Database.GetQuery("GET_HARDWARE")
	if !ok {
		return nil, 0, errors.New("запрос GET_HARDWARE не подготовлен")
	}

	count, err := r.Counter.countRecords("GET_HARDWARE_COUNT", nil)
	if err != nil {
		return nil, 0, err
	}

	rows, err := stmt.Query(offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var hardware []models.Hardware

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
			&_hardware.MgmtVlan,
			&_hardware.Description,
			&_hardware.CreatedAt,
			&_hardware.UpdatedAt,
			&_hardware.IsDelete,
			&_hardware.Node.Address.Street.Name,
			&_hardware.Node.Address.Street.Type.ShortName,
			&_hardware.Node.Address.House.Name,
			&_hardware.Node.Address.House.Type.ShortName,
			&_hardware.Type.Value,
			&_hardware.Type.TranslateValue,
			&switchName,
			&_hardware.Node.Name,
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
	fmt.Println(hardware)

	if hardware.Type.ID == 0 || hardware.Node.ID == 0 {
		return false
	}

	if hardware.Type.Value == "switch" && (hardware.Switch.ID == 0 || !hardware.IpAddress.Valid) {
		return false
	}

	return true
}
