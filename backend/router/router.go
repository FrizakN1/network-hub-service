package router

import (
	"backend/database"
	"backend/settings"
	"backend/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// Переменная для хранения конфигурации из файла settings.json
var config settings.Setting

// Initialization Функция роутинга
func Initialization(_config *settings.Setting) *gin.Engine {
	config = *_config

	router := gin.Default()

	// Middleware, проверяющий домен у отправителя запроса, если домену разрешены запросы, то выполняется запрос дальше,
	// В settings.json параметр AllowOrigin содерижт этот самый домен, которому разрешено делать запросы
	router.Use(func(c *gin.Context) {
		// Получение из settings.json разрешенной ссылки
		allowedOrigin := config.AllowOrigin

		// Получение ссылки из запроса
		origin := c.Request.Header.Get("Origin")

		// Проверка ссылок на соответсвие
		if allowedOrigin == origin {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(http.StatusOK)
			} else {
				c.Next()
			}
		} else {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(http.StatusOK)
			} else {
				c.Next()
			}
			//c.AbortWithStatus(http.StatusUnauthorized)
		}
	})

	// Групировка запросов содержащих в запросе /api в отдельный роутер routerAPI
	routerAPI := router.Group("/api")

	// Обработка запросов
	routerAPI.POST("/search", handlerGetSuggestions)
	routerAPI.GET("/get_house/:id", handlerGetHouse)
	routerAPI.POST("/upload_file", handlerUploadFile)
	routerAPI.GET("/get_files/:id", handlerGetFiles)
	routerAPI.POST("/archive_file", handlerArchiveFile)
	routerAPI.POST("/delete_file", handlerDeleteFile)
	routerAPI.POST("/get_list", handlerGetList)

	return router
}

func handlerGetList(c *gin.Context) {
	var offset int

	err := c.BindJSON(&offset)
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	addresses, count, err := database.GetList(offset)
	if err != nil {
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, gin.H{
		"Addresses": addresses,
		"Count":     count,
	})
}

func handlerDeleteFile(c *gin.Context) {
	var file database.File

	err := c.BindJSON(&file)
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	err = file.Delete()
	if err != nil {
		handlerError(c, err, 400)
		return
	}

	err = os.Remove(file.Path)
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, file)
}

func handlerArchiveFile(c *gin.Context) {
	var file database.File

	err := c.BindJSON(&file)
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	err = file.Archive()
	if err != nil {
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, file)
}

func handlerGetFiles(c *gin.Context) {
	houseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	files, err := database.GetFiles(houseID)
	if err != nil {
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, files)
}

func handlerUploadFile(c *gin.Context) {
	var uploadFile database.File
	var err error

	uploadFile.House.ID, err = strconv.Atoi(c.PostForm("houseID"))
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

	err = c.SaveUploadedFile(file, uploadFile.Path)
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	uploadFile.UploadAt = time.Now().Unix()

	err = uploadFile.Create()
	if err != nil {
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, uploadFile)
}

func handlerGetHouse(c *gin.Context) {
	var address database.Address
	var err error

	address.House.ID, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	err = address.GetHouse()
	if err != nil {
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, address)
}

func handlerGetSuggestions(c *gin.Context) {
	var search database.Search
	err := c.BindJSON(&search)
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	if search.Limit == 0 {
		search.Limit = 10
	}

	suggestions, count, err := database.GetSuggestions(search)
	if err != nil {
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, gin.H{
		"Addresses": suggestions,
		"Count":     count,
	})
}

func handlerError(c *gin.Context, err error, code int) {
	fmt.Println(err)
	c.JSON(code, nil)
	c.Abort()
}
