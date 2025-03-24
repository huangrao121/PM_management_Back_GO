package config

import (
	"log"
	"os"

	"pm_go_version/app/domain/entity"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	godotenv.Load()
}

func ConnectToDB() *gorm.DB {
	dsn := os.Getenv("DB_DSN")
	//temp := "host=localhost user=postgres password=nxbsin60w5 dbname=pmsystem port=5433 sslmode=disable"
	//fmt.Println(dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database")
	}
	db.SetupJoinTable(&entity.User{}, "Workspaces", &entity.UserWorkspace{})
	db.AutoMigrate(&entity.User{}, &entity.UserWorkspace{}, &entity.Workspace{}, &entity.Project{}, &entity.Task{})
	return db
}
