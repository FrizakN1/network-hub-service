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

func handlerGetHouseFiles(c *gin.Context) {
	houseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	files, err := database.GetHouseFiles(houseID)
	if err != nil {
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, files)
}

func handlerUploadFile(c *gin.Context) {
	var uploadFile database.File
	var err error
	var fileFor = c.PostForm("type")

	// Обработка ID (остается без изменений)
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

	// Получаем файл
	file, err := c.FormFile("file")
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	timeNow := strconv.Itoa(int(time.Now().Unix()))
	uploadFile.Name = file.Filename
	uploadFile.Path = filepath.Join("./upload", timeNow+"_"+uploadFile.Name)

	// Открываем файл для обработки
	srcFile, err := file.Open()
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}
	defer srcFile.Close()

	// Определяем тип файла по расширению
	ext := strings.ToLower(filepath.Ext(uploadFile.Name))
	isImage := false
	var img image.Image
	var format string

	// Проверяем, является ли файл изображением (jpeg/png)
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

	// Создаем файл для сохранения
	dstFile, err := os.Create(uploadFile.Path)
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}
	defer dstFile.Close()

	if isImage {
		// Оптимизируем изображение в зависимости от формата
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
		// Для не-изображений просто копируем файл как есть
		_, err = io.Copy(dstFile, srcFile)
		if err != nil {
			utils.Logger.Println(err)
			handlerError(c, err, 400)
			return
		}
	}

	uploadFile.UploadAt = time.Now().Unix()

	// Остальная логика (остается без изменений)
	if fileFor == "nodes" {
		uploadFile.IsPreviewImage, err = strconv.ParseBool(c.PostForm("onlyImage"))
		if err != nil {
			utils.Logger.Println(err)
			handlerError(c, err, 400)
			return
		}
	}

	err = uploadFile.CreateFile(strings.ToUpper(fileFor))
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, uploadFile)
}

func handlerFile(c *gin.Context) {
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
		key = "HARDWARE"
	}

	if action == "archive" {
		err = file.Archive(key)
	} else if action == "delete" {
		err = os.Remove(file.Path)
		if err != nil {
			utils.Logger.Println(err)
			handlerError(c, err, 400)
			return
		}

		err = file.Delete(key)
	}

	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, file)
}
