package router

import (
	"backend/database"
	"backend/settings"
	"backend/utils"
	"fmt"
	"github.com/dgrijalva/jwt-go"
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

var jwtSecret = []byte("!S@crEtW0r@")

type Claims struct {
	SessionHash string `json:"hash"`
	jwt.StandardClaims
}

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

		// Проверка ссылок на соответствие
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
		}
	})

	// Групировка запросов содержащих в запросе /api в отдельный роутер routerAPI
	routerAPI := router.Group("/api")

	routerAPI.POST("/login", handlerLogin)

	routerAPI.Use(authMiddleware())

	// Обработка запросов
	routerAPI.GET("/logout", handlerLogout)
	routerAPI.POST("/search", handlerGetSuggestions)
	routerAPI.GET("/get_house/:id", handlerGetHouse)
	routerAPI.POST("/upload_file", handlerUploadFile)
	routerAPI.GET("/get_house_files/:id", handlerGetHouseFiles)
	routerAPI.GET("/get_node_files/:id", handlerGetNodeFiles)
	routerAPI.GET("/get_hardware_files/:id", handlerGetHardwareFiles)
	routerAPI.GET("/get_node_images/:id", handlerGetNodeImages)
	routerAPI.POST("/archive_file", handlerArchiveFile)
	routerAPI.POST("/delete_file", handlerDeleteFile)
	routerAPI.POST("/get_list", handlerGetList)
	routerAPI.GET("/get_auth", handlerGetAuth)
	routerAPI.GET("/get_users", handlerGetUsers)
	routerAPI.POST("/change_user_status", handlerChangeUserStatus)
	routerAPI.GET("/get_roles", handlerGetRoles)
	routerAPI.POST("/create_user", handlerCreateUser)
	routerAPI.POST("/edit_user", handlerEditUser)
	routerAPI.POST("/get_nodes", handlerGetNodes)
	routerAPI.POST("/get_search_nodes", handlerGetSearchNodes)
	routerAPI.POST("/get_nodes/:id", handlerGetHouseNodes)
	routerAPI.GET("/get_node/:id", handlerGetNode)
	routerAPI.POST("/create_node", handlerCreateNode)
	routerAPI.POST("/edit_node", handlerEditNode)
	routerAPI.GET("/get_owners", handlerGetOwners)
	routerAPI.POST("/create_:reference", handlerCreateReferenceRecord)
	routerAPI.POST("/edit_:reference", handlerEditReferenceRecord)
	routerAPI.GET("/get_node_types", handlerGetNodeTypes)
	routerAPI.GET("/get_hardware_types", handlerGetHardwareTypes)
	routerAPI.GET("/get_operation_modes", handlerGetOperationModes)
	routerAPI.POST("/create_switch", handlerCreateSwitch)
	routerAPI.POST("/edit_switch", handlerEditSwitch)
	routerAPI.GET("/get_switches", handlerGetSwitches)
	routerAPI.POST("/get_hardware", handlerGetHardware)
	routerAPI.GET("/get_hardware/:id", handlerGetHardwareByID)
	routerAPI.POST("/get_search_hardware", handlerGetSearchHardware)
	routerAPI.POST("/get_house_hardware/:id", handlerGetHouseHardware)
	routerAPI.POST("/get_node_hardware/:id", handlerGetNodeHardware)
	routerAPI.POST("/create_hardware", handlerCreateHardware)
	routerAPI.POST("/edit_hardware", handlerEditHardware)

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

	var key string

	if file.House.ID > 0 {
		key = "HOUSE"
	} else if file.Node.ID > 0 {
		key = "NODE"
	} else if file.Hardware.ID > 0 {
		key = "HARDWARE"
	}

	err = file.Delete(key)
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

	var key string

	if file.House.ID > 0 {
		key = "HOUSE"
	} else if file.Node.ID > 0 {
		key = "NODE"
	} else if file.Hardware.ID > 0 {
		key = "HARDWARE"
	}

	err = file.Archive(key)
	if err != nil {
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, file)
}

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
	if fileFor == "house" {
		uploadFile.House.ID, err = strconv.Atoi(c.PostForm("id"))
	} else if fileFor == "node" {
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
	if fileFor == "node" {
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
	//var uploadFile database.File
	//var err error
	//var fileFor = c.PostForm("type")
	//
	//if fileFor == "house" {
	//	uploadFile.House.ID, err = strconv.Atoi(c.PostForm("id"))
	//	if err != nil {
	//		utils.Logger.Println(err)
	//		handlerError(c, err, 400)
	//		return
	//	}
	//} else {
	//	uploadFile.Node.ID, err = strconv.Atoi(c.PostForm("id"))
	//	if err != nil {
	//		utils.Logger.Println(err)
	//		handlerError(c, err, 400)
	//		return
	//	}
	//}
	//
	//file, err := c.FormFile("file")
	//if err != nil {
	//	utils.Logger.Println(err)
	//	handlerError(c, err, 400)
	//	return
	//}
	//
	//timeNow := strconv.Itoa(int(time.Now().Unix()))
	//uploadFile.Name = file.Filename
	//uploadFile.Path = filepath.Join("./upload", timeNow+"_"+uploadFile.Name)
	//
	//err = c.SaveUploadedFile(file, uploadFile.Path)
	//if err != nil {
	//	utils.Logger.Println(err)
	//	handlerError(c, err, 400)
	//	return
	//}
	//
	//uploadFile.UploadAt = time.Now().Unix()
	//
	//if fileFor == "house" {
	//	err = uploadFile.CreateForHouse()
	//} else {
	//	uploadFile.IsPreviewImage, err = strconv.ParseBool(c.PostForm("onlyImage"))
	//	if err != nil {
	//		utils.Logger.Println(err)
	//		handlerError(c, err, 400)
	//		return
	//	}
	//
	//	err = uploadFile.CreateForNode()
	//}
	//if err != nil {
	//	utils.Logger.Println(err)
	//	handlerError(c, err, 400)
	//	return
	//}
	//
	//c.JSON(200, uploadFile)
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

func generateToken(hash string) (string, error) {
	claims := &Claims{
		SessionHash:    hash,
		StandardClaims: jwt.StandardClaims{},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			fmt.Println("Не обнаружен заголовок авторизации")
			c.JSON(401, gin.H{"error": "Не обнаружен заголовок авторизации"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			fmt.Println("Неверный формат токена")
			c.JSON(401, gin.H{"error": "Неверный формат токена"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims := &Claims{}

		token, e := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if e != nil {
			fmt.Println(e)
			fmt.Println("Неверный токен")
			utils.Logger.Println(e)
			c.JSON(401, gin.H{"error": "Неверный токен"})
			c.Abort()
			return
		}

		if !token.Valid {
			fmt.Println("Токен не валиден")
			c.JSON(401, gin.H{"error": "Токен не валиден"})
			c.Abort()
			return
		}

		session := database.GetSession(claims.SessionHash)
		if session == nil {
			fmt.Println("Сессия не найдена")
			c.JSON(401, gin.H{"error": "Сессия не найдена"})
			c.Abort()
			return
		}

		c.Set("sessionHash", session.Hash)

		c.Next()
	}
}

func checkSession(c *gin.Context) (*database.Session, bool) {
	session, ok := c.Get("session")
	if !ok {
		return nil, false
	}

	sessionObj, ok := session.(*database.Session)
	if !ok {
		return nil, false
	}

	return sessionObj, true
}
