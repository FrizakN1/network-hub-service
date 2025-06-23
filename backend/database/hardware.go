package database

import (
	"backend/models"
	"database/sql"
	"errors"
	"github.com/lib/pq"
)

type HardwareRepository interface {
	GetHardwareByID(hardware *models.Hardware) error
	EditHardware(hardware *models.Hardware) error
	CreateHardware(hardware *models.Hardware) error
	GetHardware(offset int, houseID int, nodeID int) ([]models.Hardware, int, error)
	ValidateHardware(hardware models.Hardware) bool
	DeleteHardware(hardwareID int) error
	GetHardwareForIndex() ([]models.Hardware, error)
	GetHardwareByIDs(hardwareIDs []int32) ([]models.Hardware, error)
}

type DefaultHardwareRepository struct {
	Database Database
}

func (r *DefaultHardwareRepository) GetHardwareByIDs(hardwareIDs []int32) ([]models.Hardware, error) {
	stmt, ok := r.Database.GetQuery("GET_HARDWARE_BY_IDS")
	if !ok {
		return nil, errors.New("query GET_HARDWARE_BY_IDS is not prepare")
	}

	rows, err := stmt.Query(pq.Array(hardwareIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orderedHardware := make([]models.Hardware, 0, len(hardwareIDs))
	hardwareMap := make(map[int32]models.Hardware)

	for rows.Next() {
		var hd models.Hardware
		var switchID sql.NullInt32
		var switchName sql.NullString

		if err = rows.Scan(
			&hd.ID,
			&hd.Node.ID,
			&hd.Type.ID,
			&switchID,
			&hd.IpAddress,
			&hd.Node.HouseId,
			&hd.Type.Key,
			&hd.Type.Value,
			&switchName,
			&hd.Node.Name,
		); err != nil {
			return nil, err
		}

		if switchID.Valid {
			hd.Switch = models.Switch{ID: int(switchID.Int32), Name: switchName.String}
		}

		hardwareMap[int32(hd.ID)] = hd
	}

	for _, id := range hardwareIDs {
		if node, ok := hardwareMap[id]; ok {
			orderedHardware = append(orderedHardware, node)
		}
	}

	return orderedHardware, nil
}

func (r *DefaultHardwareRepository) GetHardwareForIndex() ([]models.Hardware, error) {
	stmt, ok := r.Database.GetQuery("GET_HARDWARE_FOR_INDEX")
	if !ok {
		return nil, errors.New("query GET_HARDWARE_FOR_INDEX is not prepare")
	}

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hardware []models.Hardware

	for rows.Next() {
		var hd models.Hardware
		var switchName sql.NullString

		if err = rows.Scan(
			&hd.ID,
			&hd.Type.Value,
			&hd.Node.Name,
			&switchName,
			&hd.IpAddress,
			&hd.Node.HouseId,
			&hd.IsDelete,
		); err != nil {
			return nil, err
		}

		if switchName.Valid {
			hd.Switch.Name = switchName.String
		}

		hardware = append(hardware, hd)
	}

	return hardware, nil
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
		&hardware.Node.HouseId,
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
			&_hardware.Node.HouseId,
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
