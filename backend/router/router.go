package router

import (
	"backend/database"
	"backend/settings"
	"backend/utils"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
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
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(http.StatusNoContent)
			} else {
				c.Next()
			}
		}
	})

	// Групировка запросов содержащих в запросе /api в отдельный роутер routerAPI
	routerAPI := router.Group("/api")

	routerAPI.POST("/auth/login", handlerLogin)

	routerAPI.Use(authMiddleware())

	routerAPI.GET("/auth/logout", handlerLogout)
	routerAPI.GET("/auth/me", handlerGetAuth)
	//routerAPI.GET("/auth/users", handlerGetUsers)

	users := routerAPI.Group("/users")
	{
		users.GET("", handlerGetUsers)
		users.POST("", handlerCreateUser)
		users.PUT("", handlerEditUser)
		users.PATCH("/status", handlerChangeUserStatus)
	}

	nodes := routerAPI.Group("/nodes")
	{
		nodes.GET("", handlerGetNodes)
		nodes.GET("/:id", handlerGetNode)
		nodes.GET("/search", handlerGetSearchNodes)
		nodes.GET("/:id/files", handlerGetNodeFiles)
		nodes.GET("/:id/images", handlerGetNodeImages)
		nodes.GET("/:id/hardware", handlerGetNodeHardware)
		nodes.POST("", handlerCreateNode)
		nodes.PUT("", handlerEditNode)
	}

	houses := routerAPI.Group("/houses")
	{
		houses.GET("", handlerGetHouses)
		houses.GET("/:id", handlerGetHouse)
		houses.GET("/search", handlerGetSuggestions)
		houses.GET("/:id/files", handlerGetHouseFiles)
		houses.GET("/:id/nodes", handlerGetHouseNodes)
		houses.GET("/:id/hardware", handlerGetHouseHardware)
	}

	hardware := routerAPI.Group("/hardware")
	{
		hardware.GET("", handlerGetHardware)
		hardware.GET("/search", handlerGetSearchHardware)
		hardware.GET("/:id", handlerGetHardwareByID)
		hardware.GET("/:id/files", handlerGetHardwareFiles)
		hardware.POST("", handlerCreateHardware)
		hardware.PUT("", handlerEditHardware)
	}

	switches := routerAPI.Group("/switches")
	{
		switches.GET("", handlerGetSwitches)
		switches.POST("", handlerCreateSwitch)
		switches.PUT("", handlerEditSwitch)
	}

	files := routerAPI.Group("/files")
	{
		files.POST("/upload", handlerUploadFile)
		files.POST("/:action", handlerFile)
	}

	references := routerAPI.Group("/references")
	{
		references.GET("/:reference", func(c *gin.Context) {
			handlerGetReference(c, false)
		})
		references.GET("/role", func(c *gin.Context) {
			handlerGetReference(c, true)
		})
		references.POST("/:reference", func(c *gin.Context) {
			handleReferenceRecord(c, false)
		})
		references.PUT("/:reference", func(c *gin.Context) {
			handleReferenceRecord(c, true)
		})
	}

	return router
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
