package main

import (
	"log"
	"os"
	"pm_go_version/app/pkg/redis_config"
	"pm_go_version/app/router"
	"pm_go_version/config"

	"github.com/joho/godotenv"

	log "github.com/sirupsen/logrus"
)

func Init() {
	godotenv.Load()
	config.InitLog()
}

func main() {
	// 初始化 Redis
	if err := redis_config.InitRedis(); err != nil {
		log.Warn("Redis failed to initialize, the app will continue to use normal mode: ", err)
	} else {
		defer redis_config.CloseRedis()
	}

	port := os.Getenv("PORT")

	init := config.Init()
	app := router.Init(init)
	app.Run(":" + port)
}
