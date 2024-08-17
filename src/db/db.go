package database

import (
	"database/sql"
	"fmt"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mdelclaro/gobrax/src/config"
	"github.com/mdelclaro/gobrax/src/repository/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Dbinstance struct {
	Db *gorm.DB
}

var DB Dbinstance

func StartDb() Dbinstance {
	if DB.Db != nil {
		return DB
	}

	// host := config.GetEnv("DB_HOST")
	host := "db"
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
		&entities.Driver{},
		&entities.Truck{},
	)

	DB = Dbinstance{
		Db: db,
	}

	return DB
}

func StartDbMock(t *testing.T) (*sql.DB, *gorm.DB, sqlmock.Sqlmock) {
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqldb,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		t.Fatal(err)
	}

	DB = Dbinstance{
		Db: db,
	}

	return sqldb, db, mock
}
