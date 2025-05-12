package main

import (
	"backend/database"
	"backend/router"
	"backend/settings"
	"backend/utils"
)

func main() {
	utils.InitLogger()

	config := settings.Load("settings.json")

	database.Connection(config)

	userService := database.NewUserService()

	if err := userService.CheckAdmin(config); err != nil {
		return
	}

	_ = router.Initialization(config).Run(config.Address + ":" + config.Port)
}
