package router

import (
	"backend/handlers"
	"backend/middleware"
	"backend/utils"
	"fmt"
	"github.com/gin-gonic/gin"
)

// Initialization Функция роутинга
func Initialization() *gin.Engine {
	userService := handlers.InitUserClient()
	logger := utils.InitLogger()

	mw := middleware.NewMiddleware(&userService, &logger)
	handlerUser := handlers.NewUserHandler(&userService)
	handlerSwitch := handlers.NewSwitchHandler()
	handlerReference := handlers.NewReferenceHandler()
	handlerNode := handlers.NewNodeHandler()
	handleHardware := handlers.NewHardwareHandler()
	handlerFile := handlers.NewFileHandler()
	handlerEvent := handlers.NewEventHandler(&userService)
	handlerAuth := handlers.NewAuthHandler(&userService)
	handlerAddress := handlers.NewAddressHandler()

	router := gin.Default()

	router.Use(mw.ErrorMiddleware())

	// Middleware, проверяющий домен у отправителя запроса, если домену разрешены запросы, то выполняется запрос дальше,
	// В settings.json параметр AllowOrigin содерижт этот самый домен, которому разрешено делать запросы
	router.Use(mw.CorsMiddleware())

	// Групировка запросов содержащих в запросе /api в отдельный роутер routerAPI
	routerAPI := router.Group("/api")

	routerAPI.POST("/auth/login", handlerAuth.HandlerLogin)

	routerAPI.Use(mw.AuthMiddleware())

	routerAPI.GET("/auth/logout", handlerAuth.HandlerLogout)
	routerAPI.GET("/auth/me", handlerAuth.HandlerGetAuth)

	users := routerAPI.Group("/users")
	{
		users.GET("", handlerUser.HandlerGetUsers)
		users.POST("", handlerUser.HandlerCreateUser)
		users.PUT("", handlerUser.HandlerEditUser)
		users.PATCH("/:id/status", handlerUser.HandlerChangeUserStatus)
	}

	nodes := routerAPI.Group("/nodes")
	{
		nodes.GET("", handlerNode.HandlerGetNodes)
		nodes.GET("/:id", handlerNode.HandlerGetNode)
		nodes.GET("/search", handlerNode.HandlerGetSearchNodes)
		nodes.GET("/:id/files", handlerFile.HandlerGetNodeFiles)
		nodes.GET("/:id/images", handlerFile.HandlerGetNodeImages)
		nodes.GET("/:id/hardware", handleHardware.HandlerGetNodeHardware)
		nodes.POST("", handlerNode.HandlerCreateNode)
		nodes.PUT("", handlerNode.HandlerEditNode)
		nodes.GET("/:id/events/:type", func(c *gin.Context) {
			handlerEvent.HandlerGetEvents(c, "NODE")
		})
		nodes.DELETE("/:id", handlerNode.HandlerDeleteNode)
	}

	houses := routerAPI.Group("/houses")
	{
		houses.GET("", handlerAddress.HandlerGetHouses)
		houses.GET("/:id", handlerAddress.HandlerGetHouse)
		houses.GET("/search", handlerAddress.HandlerGetSuggestions)
		houses.GET("/:id/files", handlerFile.HandlerGetHouseFiles)
		houses.GET("/:id/nodes", handlerNode.HandlerGetHouseNodes)
		houses.GET("/:id/hardware", handleHardware.HandlerGetHouseHardware)
		houses.GET("/:id/events/:type", func(c *gin.Context) {
			handlerEvent.HandlerGetEvents(c, "HOUSE")
		})
	}

	hardware := routerAPI.Group("/hardware")
	{
		hardware.GET("", handleHardware.HandlerGetHardware)
		hardware.GET("/search", handleHardware.HandlerGetSearchHardware)
		hardware.GET("/:id", handleHardware.HandlerGetHardwareByID)
		hardware.GET("/:id/files", handlerFile.HandlerGetHardwareFiles)
		hardware.POST("", handleHardware.HandlerCreateHardware)
		hardware.PUT("", handleHardware.HandlerEditHardware)
		hardware.GET("/:id/events/:type", func(c *gin.Context) {
			handlerEvent.HandlerGetEvents(c, "HARDWARE")
		})
		hardware.DELETE("/:id", handleHardware.HandlerDeleteHardware)
	}

	switches := routerAPI.Group("/switches")
	{
		switches.GET("", handlerSwitch.HandlerGetSwitches)
		switches.POST("", handlerSwitch.HandlerCreateSwitch)
		switches.PUT("", handlerSwitch.HandlerEditSwitch)
	}

	files := routerAPI.Group("/files")
	{
		files.POST("/upload", handlerFile.HandlerUploadFile)
		files.POST("/:action", handlerFile.HandlerFile)
	}

	references := routerAPI.Group("/references")
	{
		references.GET("/:reference", handlerReference.HandlerGetReference)
		references.POST("/:reference", func(c *gin.Context) {
			handlerReference.HandlerReferenceRecord(c, false)
		})
		references.PUT("/:reference", func(c *gin.Context) {
			handlerReference.HandlerReferenceRecord(c, true)
		})
	}

	routerAPI.GET("/events", func(c *gin.Context) {
		handlerEvent.HandlerGetEvents(c, "")
	})

	return router
}

func handlerError(c *gin.Context, err error, code int) {
	fmt.Println(err)
	c.JSON(code, nil)
	c.Abort()
}
