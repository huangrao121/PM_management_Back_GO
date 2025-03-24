package main

import (
	"os"
	"pm_go_version/app/router"
	"pm_go_version/config"

	"github.com/joho/godotenv"
)

func Init() {
	godotenv.Load()
	config.InitLog()
}

func main() {
	port := os.Getenv("PORT")

	init := config.Init()
	app := router.Init(init)
	app.Run(":" + port)
}
