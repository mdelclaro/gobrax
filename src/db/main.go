package database

import (
	"fmt"
	"log"

	"github.com/mdelclaro/gobrax/src/config"
	"github.com/mdelclaro/gobrax/src/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Dbinstance struct {
	Db *gorm.DB
}

var DB Dbinstance

func ConnectDb() {
	host := config.GetEnv("DB_HOST")
	user := config.GetEnv("DB_USER")
	pwd := config.GetEnv("DB_PASSWORD")
	dbName := config.GetEnv("DB_NAME")
	port := config.GetEnv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password='%s' dbname=%s port=%s sslmode=disable", host, user, pwd, dbName, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
	}

	db.Logger = logger.Default.LogMode(logger.Info)

	db.AutoMigrate(
		&models.Driver{},
		&models.Truck{},
	)

	DB = Dbinstance{
		Db: db,
	}
}
