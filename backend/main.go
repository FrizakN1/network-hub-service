package main

import (
	"backend/database"
	"backend/router"
	"backend/settings"
)

func main() {
	config := settings.Load("settings.json")

	database.Connection(config)

	_ = router.Initialization(config).Run(config.Address + ":" + config.Port)
}
