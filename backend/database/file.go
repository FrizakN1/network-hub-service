package database

import (
	"backend/models"
	"encoding/base64"
	"errors"
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
		return nil, errors.New("запрос GET_HARDWARE_FILES не подготовлен")
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
		return nil, errors.New("запрос GET_NODE_FILES не подготовлен")
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
		return nil, errors.New("запрос GET_HOUSE_FILES не подготовлен")
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
			&file.House.ID,
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
		return errors.New("запрос CREATE_FILE_" + fileFor + " не подготовлен")
	}

	var err error

	switch fileFor {
	case "NODES":
		err = stmt.QueryRow(
			file.Node.ID,
			file.Path,
			file.Name,
			file.UploadAt,
			file.InArchive,
			file.IsPreviewImage,
		).Scan(&file.ID)
		break
	case "HOUSES":
		err = stmt.QueryRow(
			file.House.ID,
			file.Path,
			file.Name,
			file.UploadAt,
			file.InArchive,
		).Scan(&file.ID)
		break
	case "HARDWARE":
		err = stmt.QueryRow(
			file.Hardware.ID,
			file.Path,
			file.Name,
			file.UploadAt,
			file.InArchive,
		).Scan(&file.ID)
		break
	}

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

func (r *DefaultFileRepository) Delete(file *models.File, key string) error {
	stmt, ok := r.Database.GetQuery("DELETE_FILE_" + key)
	if !ok {
		return errors.New("запрос DELETE_FILE_" + key + " не подготовлен")
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
		return errors.New("запрос ARCHIVE_FILE_" + key + " не подготовлен")
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
