package main

import (
	"log"
	"os"
	"pm_go_version/app/pkg/redis"
	"pm_go_version/app/router"
	"pm_go_version/config"

	"github.com/joho/godotenv"
)

func Init() {
	godotenv.Load()
	config.InitLog()
}

func main() {
	// 初始化 Redis
	err := redis.InitRedis(
		os.Getenv("REDIS_ADDR"),
		os.Getenv("REDIS_PASSWORD"),
		0,
	)
	if err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}

	port := os.Getenv("PORT")

	init := config.Init()
	app := router.Init(init)
	app.Run(":" + port)
}
