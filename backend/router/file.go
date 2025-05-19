package router

import (
	"backend/database"
	"backend/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type FileHandler interface {
	handlerGetHardwareFiles(c *gin.Context)
	handlerGetNodeImages(c *gin.Context)
	handlerGetNodeFiles(c *gin.Context)
	handlerGetHouseFiles(c *gin.Context)
	handlerUploadFile(c *gin.Context)
	handlerFile(c *gin.Context)
}

type DefaultFileHandler struct {
	Privilege       Privilege
	FileService     database.FileService
	EventService    database.EventService
	NodeService     database.NodeService
	HardwareService database.HardwareService
}

func NewFileHandler() FileHandler {
	return &DefaultFileHandler{
		Privilege:       &DefaultPrivilege{},
		FileService:     &database.DefaultFileService{},
		EventService:    &database.DefaultEventService{},
		NodeService:     &database.DefaultNodeService{},
		HardwareService: &database.DefaultHardwareService{},
	}
}

func (h *DefaultFileHandler) handlerGetHardwareFiles(c *gin.Context) {
	hardwareID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	files, err := h.FileService.GetHardwareFiles(hardwareID)
	if err != nil {
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, files)
}

func (h *DefaultFileHandler) handlerGetNodeImages(c *gin.Context) {
	nodeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	files, err := h.FileService.GetNodeFiles(nodeID, true)
	if err != nil {
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, files)
}

func (h *DefaultFileHandler) handlerGetNodeFiles(c *gin.Context) {
	nodeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	files, err := h.FileService.GetNodeFiles(nodeID, false)
	if err != nil {
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, files)
}

func (h *DefaultFileHandler) handlerGetHouseFiles(c *gin.Context) {
	houseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	files, err := h.FileService.GetHouseFiles(houseID)
	if err != nil {
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, files)
}

func (h *DefaultFileHandler) handlerUploadFile(c *gin.Context) {
	session, _, isOperatorOrHigher := h.Privilege.getPrivilege(c)

	if !isOperatorOrHigher {
		c.JSON(403, nil)
		return
	}
	var (
		uploadFile database.File
		err        error
		fileFor    = c.PostForm("type")
		event      database.Event
	)

	if fileFor == "houses" {
		uploadFile.House.ID, err = strconv.Atoi(c.PostForm("id"))

		event = database.Event{
			Address:  database.Address{House: database.AddressElement{ID: uploadFile.House.ID}},
			Node:     nil,
			Hardware: nil,
		}
	} else if fileFor == "nodes" {
		uploadFile.Node.ID, err = strconv.Atoi(c.PostForm("id"))

		if err = h.NodeService.GetNode(&uploadFile.Node); err != nil {
			utils.Logger.Println(err)
			handlerError(c, err, 400)
		}

		event = database.Event{
			Address:  database.Address{House: database.AddressElement{ID: uploadFile.Node.Address.House.ID}},
			Node:     &database.Node{ID: uploadFile.Node.ID},
			Hardware: nil,
		}
	} else {
		uploadFile.Hardware.ID, err = strconv.Atoi(c.PostForm("id"))

		if err = h.HardwareService.GetHardwareByID(&uploadFile.Hardware); err != nil {
			utils.Logger.Println(err)
			handlerError(c, err, 400)
		}

		event = database.Event{
			Address:  database.Address{House: database.AddressElement{ID: uploadFile.Hardware.Node.Address.House.ID}},
			Node:     &database.Node{ID: uploadFile.Hardware.Node.ID},
			Hardware: &database.Hardware{ID: uploadFile.Hardware.ID},
		}
	}
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	timeNow := strconv.Itoa(int(time.Now().Unix()))
	uploadFile.Name = file.Filename
	uploadFile.Path = filepath.Join("./upload", timeNow+"_"+uploadFile.Name)

	srcFile, err := file.Open()
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}
	defer srcFile.Close()

	ext := strings.ToLower(filepath.Ext(uploadFile.Name))
	isImage := false
	var img image.Image
	var format string

	if ext == ".jpg" || ext == ".jpeg" || ext == ".png" {
		isImage = true

		// Декодируем изображение
		img, format, err = image.Decode(srcFile)
		if err != nil {
			utils.Logger.Println(err)
			handlerError(c, err, 400)
			return
		}
	}

	dstFile, err := os.Create(uploadFile.Path)
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}
	defer dstFile.Close()

	if isImage {
		switch format {
		case "jpeg":
			// JPEG: качество 50%
			err = jpeg.Encode(dstFile, img, &jpeg.Options{Quality: 50})
		case "png":
			// PNG: максимальное сжатие
			encoder := png.Encoder{CompressionLevel: png.BestCompression}
			err = encoder.Encode(dstFile, img)
		}
		if err != nil {
			utils.Logger.Println(err)
			handlerError(c, err, 400)
			return
		}
	} else {
		_, err = io.Copy(dstFile, srcFile)
		if err != nil {
			utils.Logger.Println(err)
			handlerError(c, err, 400)
			return
		}
	}

	uploadFile.UploadAt = time.Now().Unix()

	if fileFor == "nodes" {
		uploadFile.IsPreviewImage, err = strconv.ParseBool(c.PostForm("onlyImage"))
		if err != nil {
			utils.Logger.Println(err)
			handlerError(c, err, 400)
			return
		}
	}

	err = h.FileService.CreateFile(&uploadFile, strings.ToUpper(fileFor))
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	event.UserId = session.User.Id
	event.Description = fmt.Sprintf("Загрузка файла: %s", uploadFile.Name)
	event.CreatedAt = time.Now().Unix()

	if err = h.EventService.CreateEvent(event); err != nil {
		utils.Logger.Println(err)
	}

	c.JSON(200, uploadFile)
}

func (h *DefaultFileHandler) handlerFile(c *gin.Context) {
	session, isAdmin, isOperatorOrHigher := h.Privilege.getPrivilege(c)

	action := c.Param("action")

	if action == "delete" && !isAdmin {
		c.JSON(403, nil)
		return
	}

	if action == "archive" && !isOperatorOrHigher {
		c.JSON(403, nil)
		return
	}

	var (
		key   string
		file  database.File
		event database.Event
	)

	err := c.BindJSON(&file)
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	if file.House.ID > 0 {
		key = "HOUSES"

		event = database.Event{
			Address:  database.Address{House: database.AddressElement{ID: file.House.ID}},
			Node:     nil,
			Hardware: nil,
		}
	} else if file.Node.ID > 0 {
		key = "NODES"

		if err = h.NodeService.GetNode(&file.Node); err != nil {
			utils.Logger.Println(err)
			handlerError(c, err, 400)
			return
		}

		event = database.Event{
			Address:  database.Address{House: database.AddressElement{ID: file.Node.Address.House.ID}},
			Node:     &database.Node{ID: file.Node.ID},
			Hardware: nil,
		}
	} else if file.Hardware.ID > 0 {
		key = "HARDWARE"

		if err = h.HardwareService.GetHardwareByID(&file.Hardware); err != nil {
			utils.Logger.Println(err)
			handlerError(c, err, 400)
			return
		}

		event = database.Event{
			Address:  database.Address{House: database.AddressElement{ID: file.Hardware.Node.Address.House.ID}},
			Node:     &database.Node{ID: file.Hardware.Node.ID},
			Hardware: &database.Hardware{ID: file.Hardware.ID},
		}
	}

	if action == "archive" {
		err = h.FileService.Archive(&file, key)

		if file.InArchive {
			event.Description = fmt.Sprintf("Перемещение файла %s в архив", file.Name)
		} else {
			event.Description = fmt.Sprintf("Перемещение файла %s из архива", file.Name)
		}
	} else if action == "delete" {
		event.Description = fmt.Sprintf("Удаление файла: %s", file.Name)

		err = os.Remove(file.Path)
		if err != nil {
			utils.Logger.Println(err)
			handlerError(c, err, 400)
			return
		}

		err = h.FileService.Delete(&file, key)
	}

	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	event.UserId = session.User.Id
	event.CreatedAt = time.Now().Unix()

	if err = h.EventService.CreateEvent(event); err != nil {
		utils.Logger.Println(err)
	}

	c.JSON(200, file)
}
