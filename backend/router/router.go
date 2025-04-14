package router

import (
	"backend/database"
	"backend/settings"
	"backend/utils"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
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

	routerAPI.POST("/login", handlerLogin)

	routerAPI.Use(authMiddleware())

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

		c.Set("session", session)

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
