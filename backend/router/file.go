package router

import (
	"backend/database"
	"backend/utils"
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
	FileService database.FileService
}

func (fh *DefaultFileHandler) handlerGetHardwareFiles(c *gin.Context) {
	hardwareID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	files, err := fh.FileService.GetHardwareFiles(hardwareID)
	if err != nil {
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, files)
}

func (fh *DefaultFileHandler) handlerGetNodeImages(c *gin.Context) {
	nodeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	files, err := fh.FileService.GetNodeFiles(nodeID, true)
	if err != nil {
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, files)
}

func (fh *DefaultFileHandler) handlerGetNodeFiles(c *gin.Context) {
	nodeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	files, err := fh.FileService.GetNodeFiles(nodeID, false)
	if err != nil {
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, files)
}

func (fh *DefaultFileHandler) handlerGetHouseFiles(c *gin.Context) {
	houseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	files, err := fh.FileService.GetHouseFiles(houseID)
	if err != nil {
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, files)
}

func (fh *DefaultFileHandler) handlerUploadFile(c *gin.Context) {
	var uploadFile database.File
	var err error
	var fileFor = c.PostForm("type")

	if fileFor == "houses" {
		uploadFile.House.ID, err = strconv.Atoi(c.PostForm("id"))
	} else if fileFor == "nodes" {
		uploadFile.Node.ID, err = strconv.Atoi(c.PostForm("id"))
	} else {
		uploadFile.Hardware.ID, err = strconv.Atoi(c.PostForm("id"))
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
			err = jpeg.Encode(dstFile, img, &jpeg.Options{Quality: 50})
		case "png":
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

	if fileFor == "node" {
		uploadFile.IsPreviewImage, err = strconv.ParseBool(c.PostForm("onlyImage"))
		if err != nil {
			utils.Logger.Println(err)
			handlerError(c, err, 400)
			return
		}
	}

	err = fh.FileService.CreateFile(&uploadFile, strings.ToUpper(fileFor))
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, uploadFile)
}

func (fh *DefaultFileHandler) handlerFile(c *gin.Context) {
	var file database.File

	err := c.BindJSON(&file)
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	action := c.Param("action")

	var key string

	if file.House.ID > 0 {
		key = "HOUSES"
	} else if file.Node.ID > 0 {
		key = "NODES"
	} else if file.Hardware.ID > 0 {
		key = "HARDWARES"
	}

	if action == "archive" {
		err = fh.FileService.Archive(&file, key)
	} else if action == "delete" {
		err = os.Remove(file.Path)
		if err != nil {
			utils.Logger.Println(err)
			handlerError(c, err, 400)
			return
		}

		err = fh.FileService.Delete(&file, key)
	}

	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, file)
}

func NewFileHandler() FileHandler {
	return &DefaultFileHandler{
		FileService: &database.DefaultFileService{},
	}
}
