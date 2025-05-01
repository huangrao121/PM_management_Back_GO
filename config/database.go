package config

import (
	"os"
	"pm_go_version/app/domain/entity"
	"sync"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	gdb  *gorm.DB
	once sync.Once
)

func init() {
	godotenv.Load()
}

func GetGdb() *gorm.DB {
	once.Do(func() {
		dsn := os.Getenv("DB_DSN")
		var err error
		gdb, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatal("Failed to connect to database")
		}
		gdb.SetupJoinTable(&entity.User{}, "Workspaces", &entity.UserWorkspace{})
		gdb.AutoMigrate(&entity.User{}, &entity.UserWorkspace{}, &entity.Workspace{}, &entity.Project{}, &entity.Task{})
	})
	return gdb
}

// func ConnectToDB() *gorm.DB {
// 	dsn := os.Getenv("DB_DSN")
// 	log.Info("Connecting to database with DSN: ", dsn)
// 	//temp := "host=localhost user=postgres password=nxbsin60w5 dbname=pmsystem port=5433 sslmode=disable"
// 	//fmt.Println(dsn)
// 	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
// 	if err != nil {
// 		log.Fatal("Failed to connect to database")
// 	}
// 	db.SetupJoinTable(&entity.User{}, "Workspaces", &entity.UserWorkspace{})
// 	db.AutoMigrate(&entity.User{}, &entity.UserWorkspace{}, &entity.Workspace{}, &entity.Project{}, &entity.Task{})
// 	return db
// }
