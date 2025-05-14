package database

import (
	"backend/utils"
	"database/sql"
	"encoding/base64"
	"errors"
	"io/ioutil"
)

type File struct {
	ID             int
	House          AddressElement
	Node           Node
	Hardware       Hardware
	Path           string
	Name           string
	UploadAt       int64
	Data           string
	InArchive      bool
	IsPreviewImage bool
}

type FileService interface {
	GetNodeFiles(nodeID int, onlyImage bool) ([]File, error)
	GetHouseFiles(houseID int) ([]File, error)
	CreateFile(file *File, fileFor string) error
	Delete(file *File, key string) error
	Archive(file *File, key string) error
}

type DefaultFileService struct{}

func prepareFile() []string {
	var e error
	errorsList := make([]string, 0)

	if query == nil {
		query = make(map[string]*sql.Stmt)
	}

	query["CREATE_FILE_HOUSES"], e = Link.Prepare(`
		INSERT INTO "House_files"(house_id, file_path, file_name, upload_at, in_archive) 
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["CREATE_FILE_HARDWARE"], e = Link.Prepare(`
		INSERT INTO "Hardware_files"(hardware_id, file_path, file_name, upload_at, in_archive) 
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["CREATE_FILE_NODES"], e = Link.Prepare(`
		INSERT INTO "Node_files"(node_id, file_path, file_name, upload_at, in_archive, is_preview_image) 
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["GET_HOUSE_FILES"], e = Link.Prepare(`
		SELECT * FROM "House_files" WHERE house_id = $1
		ORDER BY upload_at DESC 
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["ARCHIVE_FILE_HOUSES"], e = Link.Prepare(`
		UPDATE "House_files" SET in_archive = $2 WHERE id = $1
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["ARCHIVE_FILE_NODES"], e = Link.Prepare(`
		UPDATE "Node_files" SET in_archive = $2 WHERE id = $1
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["ARCHIVE_FILE_HARDWARE"], e = Link.Prepare(`
		UPDATE "Hardware_files" SET in_archive = $2 WHERE id = $1
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["DELETE_FILE_HOUSES"], e = Link.Prepare(`
		DELETE FROM "House_files" WHERE id = $1
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["DELETE_FILE_NODES"], e = Link.Prepare(`
		DELETE FROM "Node_files" WHERE id = $1
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["DELETE_FILE_HARDWARE"], e = Link.Prepare(`
		DELETE FROM "Hardware_files" WHERE id = $1
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	return errorsList
}

func (fs *DefaultFileService) GetNodeFiles(nodeID int, onlyImage bool) ([]File, error) {
	stmt, ok := query["GET_NODE_FILES"]
	if !ok {
		err := errors.New("запрос GET_NODE_FILES не подготовлен")
		utils.Logger.Println(err)
		return nil, err
	}

	rows, err := stmt.Query(nodeID, onlyImage)
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
			&file.Node.ID,
			&file.Path,
			&file.Name,
			&file.UploadAt,
			&file.InArchive,
			&file.IsPreviewImage,
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

func (fs *DefaultFileService) GetHouseFiles(houseID int) ([]File, error) {
	stmt, ok := query["GET_HOUSE_FILES"]
	if !ok {
		err := errors.New("запрос GET_HOUSE_FILES не подготовлен")
		utils.Logger.Println(err)
		return nil, err
	}

	rows, err := stmt.Query(houseID)
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
			&file.House.ID,
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

func (fs *DefaultFileService) CreateFile(file *File, fileFor string) error {
	stmt, ok := query["CREATE_FILE_"+fileFor]
	if !ok {
		err := "запрос CREATE_FILE_" + fileFor + " не подготовлен"
		utils.Logger.Println(err)
		return errors.New(err)
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
		utils.Logger.Println(err)
		return err
	}

	var fileData []byte

	fileData, err = ioutil.ReadFile(file.Path)
	if err != nil {
		utils.Logger.Println(err)
		return err
	}

	file.Data = base64.StdEncoding.EncodeToString(fileData)

	return nil
}

func (fs *DefaultFileService) Delete(file *File, key string) error {
	stmt, ok := query["DELETE_FILE_"+key]
	if !ok {
		err := "запрос DELETE_FILE_" + key + " не подготовлен"
		utils.Logger.Println(err)
		return errors.New(err)
	}

	_, err := stmt.Exec(file.ID)
	if err != nil {
		utils.Logger.Println(err)
		return err
	}

	return nil
}

func (fs *DefaultFileService) Archive(file *File, key string) error {
	stmt, ok := query["ARCHIVE_FILE_"+key]
	if !ok {
		err := "запрос ARCHIVE_FILE_" + key + " не подготовлен"
		utils.Logger.Println(err)
		return errors.New(err)
	}

	file.InArchive = !file.InArchive

	_, err := stmt.Exec(file.ID, file.InArchive)
	if err != nil {
		utils.Logger.Println(err)
		return err
	}

	var fileData []byte

	fileData, err = ioutil.ReadFile(file.Path)
	if err != nil {
		utils.Logger.Println(err)
		return err
	}

	file.Data = base64.StdEncoding.EncodeToString(fileData)

	return nil
}

func NewFileService() FileService {
	return &DefaultFileService{}
}
