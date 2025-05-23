package handlers

import (
	"backend/database"
	"backend/errors"
	"backend/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type FileHandler interface {
	HandlerGetHardwareFiles(c *gin.Context)
	HandlerGetNodeImages(c *gin.Context)
	HandlerGetNodeFiles(c *gin.Context)
	HandlerGetHouseFiles(c *gin.Context)
	HandlerUploadFile(c *gin.Context)
	HandlerFile(c *gin.Context)
}

type DefaultFileHandler struct {
	Privilege    Privilege
	FileRepo     database.FileRepository
	EventRepo    database.EventRepository
	NodeRepo     database.NodeRepository
	HardwareRepo database.HardwareRepository
}

func NewFileHandler(db *database.Database) FileHandler {
	return &DefaultFileHandler{
		Privilege: &DefaultPrivilege{},
		FileRepo: &database.DefaultFileRepository{
			Database: *db,
		},
		EventRepo: &database.DefaultEventRepository{
			Database: *db,
		},
		NodeRepo: &database.DefaultNodeRepository{
			Database: *db,
		},
		HardwareRepo: &database.DefaultHardwareRepository{
			Database: *db,
		},
	}
}

func (h *DefaultFileHandler) HandlerGetHardwareFiles(c *gin.Context) {
	hardwareID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to parse param(id) to int", http.StatusBadRequest))
		return
	}

	files, err := h.FileRepo.GetHardwareFiles(hardwareID)
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to get files", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, files)
}

func (h *DefaultFileHandler) HandlerGetNodeImages(c *gin.Context) {
	nodeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to parse param(id) to int", http.StatusBadRequest))
		return
	}

	files, err := h.FileRepo.GetNodeFiles(nodeID, true)
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to get files", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, files)
}

func (h *DefaultFileHandler) HandlerGetNodeFiles(c *gin.Context) {
	nodeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to parse param(id) to int", http.StatusBadRequest))
		return
	}

	files, err := h.FileRepo.GetNodeFiles(nodeID, false)
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to get files", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, files)
}

func (h *DefaultFileHandler) HandlerGetHouseFiles(c *gin.Context) {
	houseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to parse param(id) to int", http.StatusBadRequest))
		return
	}

	files, err := h.FileRepo.GetHouseFiles(houseID)
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to get files", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, files)
}

func (h *DefaultFileHandler) HandlerUploadFile(c *gin.Context) {
	session, _, isOperatorOrHigher := h.Privilege.getPrivilege(c)

	if !isOperatorOrHigher {
		c.Error(errors.NewHTTPError(nil, "forbidden", http.StatusForbidden))
		return
	}
	var (
		uploadFile models.File
		err        error
		fileFor    = c.PostForm("type")
		event      models.Event
	)

	if fileFor == "houses" {
		uploadFile.House.ID, err = strconv.Atoi(c.PostForm("id"))
		if err != nil {
			c.Error(errors.NewHTTPError(err, "failed to parse param(id) to int", http.StatusBadRequest))
			return
		}

		event = models.Event{
			Address:  models.Address{House: models.AddressElement{ID: uploadFile.House.ID}},
			Node:     nil,
			Hardware: nil,
		}
	} else if fileFor == "nodes" {
		uploadFile.Node.ID, err = strconv.Atoi(c.PostForm("id"))
		if err != nil {
			c.Error(errors.NewHTTPError(err, "failed to parse param(id) to int", http.StatusBadRequest))
			return
		}

		if err = h.NodeRepo.GetNode(&uploadFile.Node); err != nil {
			c.Error(errors.NewHTTPError(err, "failed to get node", http.StatusInternalServerError))
			return
		}

		event = models.Event{
			Address:  models.Address{House: models.AddressElement{ID: uploadFile.Node.Address.House.ID}},
			Node:     &models.Node{ID: uploadFile.Node.ID},
			Hardware: nil,
		}
	} else {
		uploadFile.Hardware.ID, err = strconv.Atoi(c.PostForm("id"))
		if err != nil {
			c.Error(errors.NewHTTPError(err, "failed to parse param(id) to int", http.StatusBadRequest))
			return
		}

		if err = h.HardwareRepo.GetHardwareByID(&uploadFile.Hardware); err != nil {
			c.Error(errors.NewHTTPError(err, "failed to get hardware", http.StatusInternalServerError))
			return
		}

		event = models.Event{
			Address:  models.Address{House: models.AddressElement{ID: uploadFile.Hardware.Node.Address.House.ID}},
			Node:     &models.Node{ID: uploadFile.Hardware.Node.ID},
			Hardware: &models.Hardware{ID: uploadFile.Hardware.ID},
		}
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to get file", http.StatusBadRequest))
		return
	}

	timeNow := strconv.Itoa(int(time.Now().Unix()))
	uploadFile.Name = file.Filename
	uploadFile.Path = filepath.Join("./upload", timeNow+"_"+uploadFile.Name)

	srcFile, err := file.Open()
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to open file", http.StatusBadRequest))
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
			c.Error(errors.NewHTTPError(err, "failed to decode image", http.StatusBadRequest))
			return
		}

		_, _ = srcFile.Seek(0, io.SeekStart)
	}

	dstFile, err := os.Create(uploadFile.Path)
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to create file", http.StatusInternalServerError))
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
			c.Error(errors.NewHTTPError(err, "failed to compress image", http.StatusInternalServerError))
			return
		}
	} else {
		_, err = io.Copy(dstFile, srcFile)
		if err != nil {
			c.Error(errors.NewHTTPError(err, "failed to copy file", http.StatusInternalServerError))
			return
		}
	}

	uploadFile.UploadAt = time.Now().Unix()

	if fileFor == "nodes" {
		uploadFile.IsPreviewImage, err = strconv.ParseBool(c.PostForm("onlyImage"))
		if err != nil {
			c.Error(errors.NewHTTPError(err, "failed to parse param(onlyImage) to bool", http.StatusBadRequest))
			return
		}
	}

	err = h.FileRepo.CreateFile(&uploadFile, strings.ToUpper(fileFor))
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to create file db", http.StatusInternalServerError))
		return
	}

	event.UserId = session.User.Id
	event.Description = fmt.Sprintf("Загрузка файла: %s", uploadFile.Name)
	event.CreatedAt = time.Now().Unix()

	if err = h.EventRepo.CreateEvent(event); err != nil {
		c.Error(errors.NewHTTPError(err, "failed to create event", http.StatusInternalServerError))
	}

	c.JSON(http.StatusOK, uploadFile)
}

func (h *DefaultFileHandler) HandlerFile(c *gin.Context) {
	session, isAdmin, isOperatorOrHigher := h.Privilege.getPrivilege(c)

	action := c.Param("action")

	if action == "delete" && !isAdmin {
		c.Error(errors.NewHTTPError(nil, "forbidden", http.StatusForbidden))
		return
	}

	if action == "archive" && !isOperatorOrHigher {
		c.Error(errors.NewHTTPError(nil, "forbidden", http.StatusForbidden))
		return
	}

	var (
		key   string
		file  models.File
		event models.Event
	)

	err := c.BindJSON(&file)
	if err != nil {
		c.Error(errors.NewHTTPError(err, "invalid json", http.StatusBadRequest))
		return
	}

	if file.House.ID > 0 {
		key = "HOUSES"

		event = models.Event{
			Address:  models.Address{House: models.AddressElement{ID: file.House.ID}},
			Node:     nil,
			Hardware: nil,
		}
	} else if file.Node.ID > 0 {
		key = "NODES"

		if err = h.NodeRepo.GetNode(&file.Node); err != nil {
			c.Error(errors.NewHTTPError(err, "failed to get node", http.StatusInternalServerError))
			return
		}

		event = models.Event{
			Address:  models.Address{House: models.AddressElement{ID: file.Node.Address.House.ID}},
			Node:     &models.Node{ID: file.Node.ID},
			Hardware: nil,
		}
	} else if file.Hardware.ID > 0 {
		key = "HARDWARE"

		if err = h.HardwareRepo.GetHardwareByID(&file.Hardware); err != nil {
			c.Error(errors.NewHTTPError(err, "failed to get hardware", http.StatusInternalServerError))
			return
		}

		event = models.Event{
			Address:  models.Address{House: models.AddressElement{ID: file.Hardware.Node.Address.House.ID}},
			Node:     &models.Node{ID: file.Hardware.Node.ID},
			Hardware: &models.Hardware{ID: file.Hardware.ID},
		}
	}

	if action == "archive" {
		err = h.FileRepo.Archive(&file, key)

		if file.InArchive {
			event.Description = fmt.Sprintf("Перемещение файла %s в архив", file.Name)
		} else {
			event.Description = fmt.Sprintf("Перемещение файла %s из архива", file.Name)
		}
	} else if action == "delete" {
		event.Description = fmt.Sprintf("Удаление файла: %s", file.Name)

		err = os.Remove(file.Path)
		if err != nil {
			c.Error(errors.NewHTTPError(err, "failed to remove file", http.StatusInternalServerError))
			return
		}

		err = h.FileRepo.Delete(&file, key)
	}

	if err != nil {
		c.Error(errors.NewHTTPError(err, fmt.Sprintf("failed to %s file", action), http.StatusInternalServerError))
		return
	}

	event.UserId = session.User.Id
	event.CreatedAt = time.Now().Unix()

	if err = h.EventRepo.CreateEvent(event); err != nil {
		c.Error(errors.NewHTTPError(err, "failed to create event", http.StatusInternalServerError))
	}

	c.JSON(http.StatusOK, file)
}
