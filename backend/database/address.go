package database

import (
	"backend/models"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
)

type AddressRepository interface {
	SetHouseParams(address *models.AddressParams) error
	GetAddressesAmounts(houseIDs []int32, offset int) (map[int32]*models.AddressParams, error)
	GetAddressParams(addressParams *models.AddressParams) error
}

type DefaultAddressRepository struct {
	addressElementTypeMap map[string]map[string]struct{}
	Database              Database
}

func (r *DefaultAddressRepository) GetAddressParams(addressParams *models.AddressParams) error {
	stmt, ok := r.Database.GetQuery("GET_ADDRESS_PARAMS")
	if !ok {
		return errors.New("query GET_ADDRESS_PARAMS is not prepare")
	}

	var (
		roofTypeID      sql.NullInt16
		wiringTypeID    sql.NullInt16
		roofTypeValue   sql.NullString
		wiringTypeValue sql.NullString
	)

	if err := stmt.QueryRow(addressParams.HouseID).Scan(
		&roofTypeID,
		&wiringTypeID,
		&roofTypeValue,
		&wiringTypeValue,
	); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if roofTypeID.Valid {
		addressParams.RoofType = models.Reference{
			ID:    int(roofTypeID.Int16),
			Value: roofTypeValue.String,
		}
	}

	if wiringTypeID.Valid {
		addressParams.WiringType = models.Reference{
			ID:    int(wiringTypeID.Int16),
			Value: wiringTypeValue.String,
		}
	}

	return nil
}

func (r *DefaultAddressRepository) GetAddressesAmounts(houseIDs []int32, offset int) (map[int32]*models.AddressParams, error) {
	key := "GET_ADDRESSES_AMOUNTS"
	var param interface{} = offset

	if houseIDs != nil {
		key = fmt.Sprintf("%s_BY_HOUSE_IDS", key)
		param = pq.Array(houseIDs)
	}

	stmt, ok := r.Database.GetQuery(key)
	if !ok {
		return nil, fmt.Errorf("query %s is not prepare\n", key)
	}

	rows, err := stmt.Query(param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	addressAmountsMap := make(map[int32]*models.AddressParams)

	for rows.Next() {
		var houseID int32
		addressAmounts := &models.AddressParams{}

		if err = rows.Scan(
			&houseID,
			&addressAmounts.FileAmount,
			&addressAmounts.NodeAmount,
			&addressAmounts.HardwareAmount,
		); err != nil {
			return nil, err
		}

		addressAmountsMap[houseID] = addressAmounts
	}

	return addressAmountsMap, nil
}

func (r *DefaultAddressRepository) SetHouseParams(address *models.AddressParams) error {
	stmt, ok := r.Database.GetQuery("SET_HOUSE_PARAMS")
	if !ok {
		return errors.New("query SET_HOUSE_PARAMS is not prepare")
	}

	_, err := stmt.Exec(address.HouseID, address.RoofType.ID, address.WiringType.ID)
	if err != nil {
		return err
	}

	return nil
}
