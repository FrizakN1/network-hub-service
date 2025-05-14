package main

import (
	"backend/database"
	"backend/router"
	"backend/utils"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	utils.InitLogger()

	if err := godotenv.Load(); err != nil {
		log.Fatalln(err)
		return
	}

	database.Connection()

	_ = router.Initialization().Run(fmt.Sprintf("%s:%s", os.Getenv("APP_ADDRESS"), os.Getenv("APP_PORT")))
}
