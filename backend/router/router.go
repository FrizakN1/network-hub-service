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

	userHandler := NewUserHandler()
	switchHandler := NewSwitchHandler()
	referenceHandler := NewReferenceHandler()
	nodeHandler := NewNodeHandler()
	addressHandler := NewAddressHandler()
	hardwareHandler := NewHardwareHandler()
	fileHandler := NewFileHandler()

	InitUserClient()

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

	routerAPI.POST("/auth/login", userHandler.handlerLogin)

	routerAPI.Use(authMiddleware())

	routerAPI.GET("/auth/logout", userHandler.handlerLogout)
	routerAPI.GET("/auth/me", userHandler.handlerGetAuth)
	//routerAPI.GET("/auth/users", handlerGetUsers)

	users := routerAPI.Group("/users")
	{
		//users.GET("", userHandler.handlerGetUsers)
		users.GET("", handlerGetUsers)
		users.POST("", userHandler.handlerCreateUser)
		users.PUT("", userHandler.handlerEditUser)
		users.PATCH("/status", userHandler.handlerChangeUserStatus)
	}

	nodes := routerAPI.Group("/nodes")
	{
		nodes.GET("", nodeHandler.handlerGetNodes)
		nodes.GET("/:id", nodeHandler.handlerGetNode)
		nodes.GET("/search", nodeHandler.handlerGetSearchNodes)
		nodes.GET("/:id/files", fileHandler.handlerGetNodeFiles)
		nodes.GET("/:id/images", fileHandler.handlerGetNodeImages)
		nodes.GET("/:id/hardware", hardwareHandler.handlerGetNodeHardware)
		nodes.POST("", nodeHandler.handlerCreateNode)
		nodes.PUT("", nodeHandler.handlerEditNode)
		nodes.GET("/:id/events/:type", func(c *gin.Context) {
			handlerGetEventsFrom(c, "NODE")
		})
		nodes.DELETE("/:id", nodeHandler.handlerDeleteNode)
	}

	houses := routerAPI.Group("/houses")
	{
		houses.GET("", addressHandler.handlerGetHouses)
		houses.GET("/:id", addressHandler.handlerGetHouse)
		houses.GET("/search", addressHandler.handlerGetSuggestions)
		houses.GET("/:id/files", fileHandler.handlerGetHouseFiles)
		houses.GET("/:id/nodes", nodeHandler.handlerGetHouseNodes)
		houses.GET("/:id/hardware", hardwareHandler.handlerGetHouseHardware)
		houses.GET("/:id/events/:type", func(c *gin.Context) {
			handlerGetEventsFrom(c, "HOUSE")
		})
	}

	hardware := routerAPI.Group("/hardware")
	{
		hardware.GET("", hardwareHandler.handlerGetHardware)
		hardware.GET("/search", hardwareHandler.handlerGetSearchHardware)
		hardware.GET("/:id", hardwareHandler.handlerGetHardwareByID)
		hardware.GET("/:id/files", fileHandler.handlerGetHardwareFiles)
		hardware.POST("", hardwareHandler.handlerCreateHardware)
		hardware.PUT("", hardwareHandler.handlerEditHardware)
		hardware.GET("/:id/events/:type", func(c *gin.Context) {
			handlerGetEventsFrom(c, "HARDWARE")
		})
		hardware.DELETE("/:id", hardwareHandler.handlerDeleteHardware)
	}

	switches := routerAPI.Group("/switches")
	{
		switches.GET("", switchHandler.handlerGetSwitches)
		switches.POST("", switchHandler.handlerCreateSwitch)
		switches.PUT("", switchHandler.handlerEditSwitch)
	}

	files := routerAPI.Group("/files")
	{
		files.POST("/upload", fileHandler.handlerUploadFile)
		files.POST("/:action", fileHandler.handlerFile)
	}

	references := routerAPI.Group("/references")
	{
		references.GET("/:reference", func(c *gin.Context) {
			referenceHandler.handlerGetReference(c, false)
		})
		references.GET("/role", func(c *gin.Context) {
			referenceHandler.handlerGetReference(c, true)
		})
		references.POST("/:reference", func(c *gin.Context) {
			referenceHandler.handleReferenceRecord(c, false)
		})
		references.PUT("/:reference", func(c *gin.Context) {
			referenceHandler.handleReferenceRecord(c, true)
		})
	}

	routerAPI.GET("/events", handlerGetEvents)

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
