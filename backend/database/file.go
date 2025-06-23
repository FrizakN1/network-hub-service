package database

import (
	"backend/models"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
)

type FileRepository interface {
	GetHardwareFiles(hardwareID int) ([]models.File, error)
	GetNodeFiles(nodeID int, onlyImage bool) ([]models.File, error)
	GetHouseFiles(houseID int) ([]models.File, error)
	CreateFile(file *models.File, fileFor string) error
	Delete(file *models.File, key string) error
	Archive(file *models.File, key string) error
}

type DefaultFileRepository struct {
	Database Database
}

func (r *DefaultFileRepository) GetHardwareFiles(hardwareID int) ([]models.File, error) {
	stmt, ok := r.Database.GetQuery("GET_HARDWARE_FILES")
	if !ok {
		return nil, errors.New("query GET_HARDWARE_FILES is not prepare")
	}

	rows, err := stmt.Query(hardwareID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []models.File
	for rows.Next() {
		var file models.File

		err = rows.Scan(
			&file.ID,
			&file.Hardware.ID,
			&file.Path,
			&file.Name,
			&file.UploadAt,
			&file.InArchive,
		)
		if err != nil {
			return nil, err
		}

		var fileData []byte

		fileData, err = ioutil.ReadFile(file.Path)
		if err != nil {
			return nil, err
		}

		file.Data = base64.StdEncoding.EncodeToString(fileData)

		files = append(files, file)
	}

	return files, nil
}

func (r *DefaultFileRepository) GetNodeFiles(nodeID int, onlyImage bool) ([]models.File, error) {
	stmt, ok := r.Database.GetQuery("GET_NODE_FILES")
	if !ok {
		return nil, errors.New("query GET_NODE_FILES is not prepare")
	}

	rows, err := stmt.Query(nodeID, onlyImage)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []models.File
	for rows.Next() {
		var file models.File

		err = rows.Scan(
			&file.ID,
			&file.Node.ID,
			&file.Path,
			&file.Name,
			&file.UploadAt,
			&file.InArchive,
			&file.IsPreviewImage,
		)
		if err != nil {
			return nil, err
		}

		var fileData []byte

		fileData, err = ioutil.ReadFile(file.Path)
		if err != nil {
			return nil, err
		}

		file.Data = base64.StdEncoding.EncodeToString(fileData)

		files = append(files, file)
	}

	return files, nil
}

func (r *DefaultFileRepository) GetHouseFiles(houseID int) ([]models.File, error) {
	stmt, ok := r.Database.GetQuery("GET_HOUSE_FILES")
	if !ok {
		return nil, errors.New("query GET_HOUSE_FILES is not prepare")
	}

	rows, err := stmt.Query(houseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []models.File
	for rows.Next() {
		var file models.File

		err = rows.Scan(
			&file.ID,
			&file.HouseId,
			&file.Path,
			&file.Name,
			&file.UploadAt,
			&file.InArchive,
		)
		if err != nil {
			return nil, err
		}

		var fileData []byte

		fileData, err = ioutil.ReadFile(file.Path)
		if err != nil {
			return nil, err
		}

		file.Data = base64.StdEncoding.EncodeToString(fileData)

		files = append(files, file)
	}

	return files, nil
}

func (r *DefaultFileRepository) CreateFile(file *models.File, fileFor string) error {
	stmt, ok := r.Database.GetQuery("CREATE_FILE_" + fileFor)
	if !ok {
		return errors.New("query CREATE_FILE_" + fileFor + " is not prepare")
	}

	var params []interface{}

	switch fileFor {
	case "NODES":
		params = []interface{}{
			file.Node.ID,
			file.Path,
			file.Name,
			file.UploadAt,
			file.InArchive,
			file.IsPreviewImage,
		}
	case "HOUSES":
		params = []interface{}{
			file.HouseId,
			file.Path,
			file.Name,
			file.UploadAt,
			file.InArchive,
		}
	case "HARDWARE":
		params = []interface{}{
			file.Hardware.ID,
			file.Path,
			file.Name,
			file.UploadAt,
			file.InArchive,
		}
	default:
		return fmt.Errorf("type is unsupported (%s)", fileFor)
	}

	if err := stmt.QueryRow(params...).Scan(&file.ID); err != nil {
		return err
	}

	var fileData []byte

	fileData, err := ioutil.ReadFile(file.Path)
	if err != nil {
		return err
	}

	file.Data = base64.StdEncoding.EncodeToString(fileData)

	return nil
}

func (r *DefaultFileRepository) Delete(file *models.File, key string) error {
	stmt, ok := r.Database.GetQuery("DELETE_FILE_" + key)
	if !ok {
		return errors.New("query DELETE_FILE_" + key + " is not prepare")
	}

	_, err := stmt.Exec(file.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *DefaultFileRepository) Archive(file *models.File, key string) error {
	stmt, ok := r.Database.GetQuery("ARCHIVE_FILE_" + key)
	if !ok {
		return errors.New("query ARCHIVE_FILE_" + key + " is not prepare")
	}

	file.InArchive = !file.InArchive

	_, err := stmt.Exec(file.ID, file.InArchive)
	if err != nil {
		return err
	}

	var fileData []byte

	fileData, err = ioutil.ReadFile(file.Path)
	if err != nil {
		return err
	}

	file.Data = base64.StdEncoding.EncodeToString(fileData)

	return nil
}
