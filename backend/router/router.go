package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

// Initialization Функция роутинга
func Initialization() *gin.Engine {
	//handler := NewHandler()

	userService := InitUserClient()

	handlerUser := NewUserHandler(&userService)
	handlerSwitch := NewSwitchHandler()
	handlerReference := NewReferenceHandler()
	handlerNode := NewNodeHandler()
	handleHardware := NewHardwareHandler()
	handlerFile := NewFileHandler()
	handlerEvent := NewEventHandler(&userService)
	handlerAuth := NewAuthHandler(&userService)
	handlerAddress := NewAddressHandler()

	router := gin.Default()

	// Middleware, проверяющий домен у отправителя запроса, если домену разрешены запросы, то выполняется запрос дальше,
	// В settings.json параметр AllowOrigin содерижт этот самый домен, которому разрешено делать запросы
	router.Use(func(c *gin.Context) {
		// Получение из settings.json разрешенной ссылки
		allowedOrigin := os.Getenv("ALLOW_ORIGIN")

		// Получение ссылки из запроса
		origin := c.Request.Header.Get("Origin")

		// Проверка ссылок на соответствие
		if allowedOrigin == origin {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(http.StatusNoContent)
				return
			} else {
				c.Next()
			}
		}
	})

	// Групировка запросов содержащих в запросе /api в отдельный роутер routerAPI
	routerAPI := router.Group("/api")

	routerAPI.POST("/auth/login", handlerAuth.handlerLogin)

	routerAPI.Use(handlerAuth.authMiddleware())

	routerAPI.GET("/auth/logout", handlerAuth.handlerLogout)
	routerAPI.GET("/auth/me", handlerAuth.handlerGetAuth)

	users := routerAPI.Group("/users")
	{
		users.GET("", handlerUser.handlerGetUsers)
		users.POST("", handlerUser.handlerCreateUser)
		users.PUT("", handlerUser.handlerEditUser)
		users.PATCH("/:id/status", handlerUser.handlerChangeUserStatus)
	}

	nodes := routerAPI.Group("/nodes")
	{
		nodes.GET("", handlerNode.handlerGetNodes)
		nodes.GET("/:id", handlerNode.handlerGetNode)
		nodes.GET("/search", handlerNode.handlerGetSearchNodes)
		nodes.GET("/:id/files", handlerFile.handlerGetNodeFiles)
		nodes.GET("/:id/images", handlerFile.handlerGetNodeImages)
		nodes.GET("/:id/hardware", handleHardware.handlerGetNodeHardware)
		nodes.POST("", handlerNode.handlerCreateNode)
		nodes.PUT("", handlerNode.handlerEditNode)
		nodes.GET("/:id/events/:type", func(c *gin.Context) {
			handlerEvent.handlerGetEvents(c, "NODE")
		})
		nodes.DELETE("/:id", handlerNode.handlerDeleteNode)
	}

	houses := routerAPI.Group("/houses")
	{
		houses.GET("", handlerAddress.handlerGetHouses)
		houses.GET("/:id", handlerAddress.handlerGetHouse)
		houses.GET("/search", handlerAddress.handlerGetSuggestions)
		houses.GET("/:id/files", handlerFile.handlerGetHouseFiles)
		houses.GET("/:id/nodes", handlerNode.handlerGetHouseNodes)
		houses.GET("/:id/hardware", handleHardware.handlerGetHouseHardware)
		houses.GET("/:id/events/:type", func(c *gin.Context) {
			handlerEvent.handlerGetEvents(c, "HOUSE")
		})
	}

	hardware := routerAPI.Group("/hardware")
	{
		hardware.GET("", handleHardware.handlerGetHardware)
		hardware.GET("/search", handleHardware.handlerGetSearchHardware)
		hardware.GET("/:id", handleHardware.handlerGetHardwareByID)
		hardware.GET("/:id/files", handlerFile.handlerGetHardwareFiles)
		hardware.POST("", handleHardware.handlerCreateHardware)
		hardware.PUT("", handleHardware.handlerEditHardware)
		hardware.GET("/:id/events/:type", func(c *gin.Context) {
			handlerEvent.handlerGetEvents(c, "HARDWARE")
		})
		hardware.DELETE("/:id", handleHardware.handlerDeleteHardware)
	}

	switches := routerAPI.Group("/switches")
	{
		switches.GET("", handlerSwitch.handlerGetSwitches)
		switches.POST("", handlerSwitch.handlerCreateSwitch)
		switches.PUT("", handlerSwitch.handlerEditSwitch)
	}

	files := routerAPI.Group("/files")
	{
		files.POST("/upload", handlerFile.handlerUploadFile)
		files.POST("/:action", handlerFile.handlerFile)
	}

	references := routerAPI.Group("/references")
	{
		references.GET("/:reference", handlerReference.handlerGetReference)
		references.POST("/:reference", func(c *gin.Context) {
			handlerReference.handleReferenceRecord(c, false)
		})
		references.PUT("/:reference", func(c *gin.Context) {
			handlerReference.handleReferenceRecord(c, true)
		})
	}

	routerAPI.GET("/events", func(c *gin.Context) {
		handlerEvent.handlerGetEvents(c, "")
	})

	return router
}

func handlerError(c *gin.Context, err error, code int) {
	fmt.Println(err)
	c.JSON(code, nil)
	c.Abort()
}
