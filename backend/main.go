package main

import (
	"backend/database"
	"backend/router"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	migration := flag.String("migration", "up", "Migration direction: up or down")
	flag.Parse()

	// Загружаем переменные среды
	if err := godotenv.Load(); err != nil {
		log.Fatalln(err)
		return
	}

	// Инициализируем базу данных
	db, err := database.InitDatabase(*migration)
	if err != nil {
		log.Fatalln(err)
		return
	}

	// Инициализируем роутер: передаем указатель на БД, адрес и порт берем из переменных среды
	_ = router.Initialization(&db).Run(fmt.Sprintf("%s:%s", os.Getenv("APP_ADDRESS"), os.Getenv("APP_PORT")))
}
