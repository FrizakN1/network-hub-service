package router

import (
	"backend/database"
	"backend/proto/userpb"
	"backend/utils"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strings"
)

type Handler interface {
	handlerGetEventsFrom(c *gin.Context, from string)
	handlerGetEvents(c *gin.Context)
	handlerGetHardwareFiles(c *gin.Context)
	handlerGetNodeImages(c *gin.Context)
	handlerGetNodeFiles(c *gin.Context)
	handlerGetHouseFiles(c *gin.Context)
	handlerUploadFile(c *gin.Context)
	handlerFile(c *gin.Context)
	handlerGetHardwareByID(c *gin.Context)
	handlerEditHardware(c *gin.Context)
	handlerCreateHardware(c *gin.Context)
	handlerGetSearchHardware(c *gin.Context)
	handlerGetNodeHardware(c *gin.Context)
	handlerGetHouseHardware(c *gin.Context)
	handlerGetHardware(c *gin.Context)
	handlerDeleteHardware(c *gin.Context)
	handlerGetHouses(c *gin.Context)
	handlerGetHouse(c *gin.Context)
	handlerGetSuggestions(c *gin.Context)
	handlerGetSearchNodes(c *gin.Context)
	handlerEditNode(c *gin.Context)
	handlerCreateNode(c *gin.Context)
	handlerGetNode(c *gin.Context)
	handlerGetHouseNodes(c *gin.Context)
	handlerGetNodes(c *gin.Context)
	handlerDeleteNode(c *gin.Context)
	handleReferenceRecord(c *gin.Context, isEdit bool)
	handlerGetReference(c *gin.Context)
	handlerGetSwitches(c *gin.Context)
	handlerEditSwitch(c *gin.Context)
	handlerCreateSwitch(c *gin.Context)
	handlerEditUser(c *gin.Context)
	handlerCreateUser(c *gin.Context)
	handlerChangeUserStatus(c *gin.Context)
	handlerGetUsers(c *gin.Context)
	handlerGetAuth(c *gin.Context)
	handlerLogout(c *gin.Context)
	handlerLogin(c *gin.Context)
}

type DefaultHandler struct {
	AddressService   database.AddressService
	SwitchService    database.SwitchService
	ReferenceService database.ReferenceService
	NodeService      database.NodeService
	HardwareService  database.HardwareService
	FileService      database.FileService
	EventService     database.EventService
}

// Initialization Функция роутинга
func Initialization() *gin.Engine {

	handler := NewHandler()

	InitUserClient()

	router := gin.Default()

	// Middleware, проверяющий домен у отправителя запроса, если домену разрешены запросы, то выполняется запрос дальше,
	// В settings.json параметр AllowOrigin содерижт этот самый домен, которому разрешено делать запросы
	router.Use(func(c *gin.Context) {
		// Получение из settings.json разрешенной ссылки
		allowedOrigin := os.Getenv("ALLOWED_ORIGIN")

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

	routerAPI.POST("/auth/login", handler.handlerLogin)

	routerAPI.Use(authMiddleware())

	routerAPI.GET("/auth/logout", handler.handlerLogout)
	routerAPI.GET("/auth/me", handler.handlerGetAuth)
	//routerAPI.GET("/auth/users", handlerGetUsers)

	users := routerAPI.Group("/users")
	{
		//users.GET("", userHandler.handlerGetUsers)
		users.GET("", handler.handlerGetUsers)
		//users.POST("", userHandler.handlerCreateUser)
		users.POST("", handler.handlerCreateUser)
		users.PUT("", handler.handlerEditUser)
		users.PATCH("/:id/status", handler.handlerChangeUserStatus)
	}

	nodes := routerAPI.Group("/nodes")
	{
		nodes.GET("", handler.handlerGetNodes)
		nodes.GET("/:id", handler.handlerGetNode)
		nodes.GET("/search", handler.handlerGetSearchNodes)
		nodes.GET("/:id/files", handler.handlerGetNodeFiles)
		nodes.GET("/:id/images", handler.handlerGetNodeImages)
		nodes.GET("/:id/hardware", handler.handlerGetNodeHardware)
		nodes.POST("", handler.handlerCreateNode)
		nodes.PUT("", handler.handlerEditNode)
		nodes.GET("/:id/events/:type", func(c *gin.Context) {
			handler.handlerGetEventsFrom(c, "NODE")
		})
		nodes.DELETE("/:id", handler.handlerDeleteNode)
	}

	houses := routerAPI.Group("/houses")
	{
		houses.GET("", handler.handlerGetHouses)
		houses.GET("/:id", handler.handlerGetHouse)
		houses.GET("/search", handler.handlerGetSuggestions)
		houses.GET("/:id/files", handler.handlerGetHouseFiles)
		houses.GET("/:id/nodes", handler.handlerGetHouseNodes)
		houses.GET("/:id/hardware", handler.handlerGetHouseHardware)
		houses.GET("/:id/events/:type", func(c *gin.Context) {
			handler.handlerGetEventsFrom(c, "HOUSE")
		})
	}

	hardware := routerAPI.Group("/hardware")
	{
		hardware.GET("", handler.handlerGetHardware)
		hardware.GET("/search", handler.handlerGetSearchHardware)
		hardware.GET("/:id", handler.handlerGetHardwareByID)
		hardware.GET("/:id/files", handler.handlerGetHardwareFiles)
		hardware.POST("", handler.handlerCreateHardware)
		hardware.PUT("", handler.handlerEditHardware)
		hardware.GET("/:id/events/:type", func(c *gin.Context) {
			handler.handlerGetEventsFrom(c, "HARDWARE")
		})
		hardware.DELETE("/:id", handler.handlerDeleteHardware)
	}

	switches := routerAPI.Group("/switches")
	{
		switches.GET("", handler.handlerGetSwitches)
		switches.POST("", handler.handlerCreateSwitch)
		switches.PUT("", handler.handlerEditSwitch)
	}

	files := routerAPI.Group("/files")
	{
		files.POST("/upload", handler.handlerUploadFile)
		files.POST("/:action", handler.handlerFile)
	}

	references := routerAPI.Group("/references")
	{
		references.GET("/:reference", handler.handlerGetReference)
		references.GET("/role", handler.handlerGetReference)
		references.POST("/:reference", func(c *gin.Context) {
			handler.handleReferenceRecord(c, false)
		})
		references.PUT("/:reference", func(c *gin.Context) {
			handler.handleReferenceRecord(c, true)
		})
	}

	routerAPI.GET("/events", handler.handlerGetEvents)

	return router
}

func handlerError(c *gin.Context, err error, code int) {
	fmt.Println(err)
	c.JSON(code, nil)
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			fmt.Println("Не обнаружен заголовок авторизации")
			c.JSON(401, gin.H{"error": "Не обнаружен заголовок авторизации"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			fmt.Println("Неверный формат токена")
			c.JSON(401, gin.H{"error": "Неверный формат токена"})
			return
		}

		tokenString := parts[1]

		session, err := userClient.GetSession(context.Background(), &userpb.GetSessionRequest{Hash: tokenString})
		if err != nil {
			utils.Logger.Println(err)
			c.JSON(400, gin.H{"error": "ошибка при получении сессии"})
			return
		}

		if session == nil {
			fmt.Println("Сессия не найдена")
			c.JSON(401, gin.H{"error": "Сессия не найдена"})
			return
		}

		c.Set("session", session)

		c.Next()
	}
}

func NewHandler() Handler {
	return &DefaultHandler{
		AddressService:   &database.DefaultAddressService{},
		SwitchService:    &database.DefaultSwitchService{},
		ReferenceService: &database.DefaultReferenceService{},
		NodeService:      &database.DefaultNodeService{},
		HardwareService:  &database.DefaultHardwareService{},
		FileService:      &database.DefaultFileService{},
		EventService:     &database.DefaultEventService{},
	}
}
