package router

import (
	"backend/database"
	"backend/handlers"
	"backend/kafka"
	"backend/middleware"
	"backend/utils"
	"context"
	"github.com/gin-gonic/gin"
	"log"
)

// Initialization Функция инициализации роутинга
func Initialization(db *database.Database) *gin.Engine {
	userService := handlers.InitUserClient() // Инициализируем связь по gRPC с user-service
	addressService := handlers.InitAddressClient()
	searchNodeService := handlers.InitSearchService()
	logger := utils.InitLogger() // Инициализируем logger

	mw := middleware.NewMiddleware(userService, &logger) // Инициализируем все middleware
	// Инициализируем хендлеры
	handlerUser := handlers.NewUserHandler(userService)
	handlerSwitch := handlers.NewSwitchHandler(db)
	handlerReference := handlers.NewReferenceHandler(db)
	handlerNode := handlers.NewNodeHandler(addressService, searchNodeService, db, &logger)
	handlerHardware := handlers.NewHardwareHandler(addressService, searchNodeService, db, &logger)
	handlerFile := handlers.NewFileHandler(db)
	handlerEvent := handlers.NewEventHandler(userService, addressService, db)
	handlerAuth := handlers.NewAuthHandler(userService)
	handlerAddress := handlers.NewAddressHandler(addressService, db)
	handlerReport := handlers.NewReportHandler(db, &logger)

	go func() {
		if err := kafka.CreateTopics(); err != nil {
			log.Fatalln(err)
		}

		if err := handlerNode.SendBatchNodes(context.Background()); err != nil {
			log.Println(err)
		}

		if err := handlerHardware.SendBatchHardware(context.Background()); err != nil {
			log.Println(err)
		}
	}()

	router := gin.Default() // Инициализируем роутер

	router.Use(mw.ErrorMiddleware()) // Говорим роутеру использовать ErrorMiddleware перед запросами для обработки ошибок возникших в запросах
	router.Use(mw.CorsMiddleware())  // Говорим роутеру использовать CorsMiddleware перед запросами для настройки CORS политики

	routerAPI := router.Group("/api")

	routerAPI.POST("/auth/login", handlerAuth.HandlerLogin)

	routerAPI.Use(mw.AuthMiddleware()) // Говорим роутеру использовать AuthMiddleware перед запросами ниже для аутентификации пользователя

	routerAPI.GET("/auth/logout", handlerAuth.HandlerLogout)
	routerAPI.GET("/auth/me", handlerAuth.HandlerGetAuth)

	users := routerAPI.Group("/users")
	{
		users.GET("", handlerUser.HandlerGetUsers)
		users.POST("", handlerUser.HandlerCreateUser)
		users.PUT("", handlerUser.HandlerEditUser)
		users.PATCH("/:id/status", handlerUser.HandlerChangeUserStatus)
		users.GET("/roles", handlerUser.GetRoles)
	}

	nodes := routerAPI.Group("/nodes")
	{
		nodes.GET("", handlerNode.HandlerGetNodes)
		nodes.GET("/:id", handlerNode.HandlerGetNode)
		nodes.GET("/search", handlerNode.HandlerGetSearchNodes)
		nodes.GET("/:id/files", handlerFile.HandlerGetNodeFiles)
		nodes.GET("/:id/images", handlerFile.HandlerGetNodeImages)
		nodes.GET("/:id/hardware", handlerHardware.HandlerGetNodeHardware)
		nodes.POST("", handlerNode.HandlerCreateNode)
		nodes.PUT("", handlerNode.HandlerEditNode)
		nodes.GET("/:id/events/:type", func(c *gin.Context) {
			handlerEvent.HandlerGetEvents(c, "NODE")
		})
		nodes.DELETE("/:id", handlerNode.HandlerDeleteNode)
		//nodes.GET("/index", handlerNode.HandlerIndexNodes)
	}

	houses := routerAPI.Group("/houses")
	{
		houses.GET("", handlerAddress.HandlerGetHouses)
		houses.GET("/:id", handlerAddress.HandlerGetHouse)
		houses.GET("/search", handlerAddress.HandlerSearchAddresses)
		houses.GET("/:id/files", handlerFile.HandlerGetHouseFiles)
		houses.GET("/:id/nodes", handlerNode.HandlerGetHouseNodes)
		houses.GET("/:id/hardware", handlerHardware.HandlerGetHouseHardware)
		houses.GET("/:id/events/:type", func(c *gin.Context) {
			handlerEvent.HandlerGetEvents(c, "HOUSE")
		})
		houses.POST("/:id/params", handlerAddress.HandlerSetHouseParams)
		houses.GET("/:id/excel", handlerNode.HandlerGetNodesExcel)
	}

	hardware := routerAPI.Group("/hardware")
	{
		hardware.GET("", handlerHardware.HandlerGetHardware)
		hardware.GET("/search", handlerHardware.HandlerGetSearchHardware)
		hardware.GET("/:id", handlerHardware.HandlerGetHardwareByID)
		hardware.GET("/:id/files", handlerFile.HandlerGetHardwareFiles)
		hardware.POST("", handlerHardware.HandlerCreateHardware)
		hardware.PUT("", handlerHardware.HandlerEditHardware)
		hardware.GET("/:id/events/:type", func(c *gin.Context) {
			handlerEvent.HandlerGetEvents(c, "HARDWARE")
		})
		hardware.DELETE("/:id", handlerHardware.HandlerDeleteHardware)
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

	report := routerAPI.Group("/report")
	{
		report.GET("", handlerReport.HandlerGetReportData)
		report.PUT("", handlerReport.HandlerEditReportData)
	}

	routerAPI.GET("/events", func(c *gin.Context) {
		handlerEvent.HandlerGetEvents(c, "")
	})

	return router
}
