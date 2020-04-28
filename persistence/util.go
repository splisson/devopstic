package persistence

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/prometheus/common/log"
	"github.com/splisson/devopstic/entities"
	"os"
)

func CreateTables(db *gorm.DB) {
	db.CreateTable(&entities.User{})
	db.CreateTable(&entities.Event{})
	db.CreateTable(&entities.Commit{})
	db.CreateTable(&entities.Incident{})
	db.CreateTable(&entities.Deployment{})
	db.AutoMigrate(&entities.Deployment{})
	db.AutoMigrate(&entities.User{})
	db.AutoMigrate(&entities.Event{})
	db.AutoMigrate(&entities.Commit{})
	db.AutoMigrate(&entities.Incident{})
}

func NewSQLiteConnection(dbId string) (*gorm.DB, string) {
	filename := fmt.Sprintf("/tmp/opstic_test_%s.db", dbId)
	db, err := gorm.Open("sqlite3", filename)
	if err != nil {
		log.Errorf("Error while connecting to sqlite local db %v", err)
		panic("failed to connect to database")
	}
	return db, filename
}

func CreatePostgresDBConnection(url string, port string, username string, password string, databaseName string) *gorm.DB {
	log.Infof("Connection to Postgres : host=%s port=%s user=%s password=XXX dbname=%s sslmode=disable", url, port, username, databaseName)
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", url, port, username, password, databaseName)
	db, err := gorm.Open("postgres", connectionString)
	if err != nil {
		log.Error(err)
		panic("failed to connect to database")
	}
	return db
}

func NewPostgresqlConnectionLocalhost() *gorm.DB {
	return CreatePostgresDBConnection("localhost", "5432", os.Getenv("DATABASE_USERNAME"), os.Getenv("DATABASE_PASSWORD"), os.Getenv("DATABASE_NAME"))
}

func NewPostgresqlConnectionWithEnv() *gorm.DB {
	url := os.Getenv("DEVOPSTIC_DATABASE_HOST")
	if url == "" {
		url = "localhost"
	}
	port := os.Getenv("DEVOPSTIC_DATABASE_PORT")
	if port == "" {
		port = "5432"
	}
	username := os.Getenv("DEVOPSTIC_DATABASE_USERNAME")
	password := os.Getenv("DEVOPSTIC_DATABASE_PASSWORD")
	databaseName := os.Getenv("DEVOPSTIC_DATABASE_NAME")
	if databaseName == "" {
		databaseName = "devopstic"
	}
	return CreatePostgresDBConnection(url, port, username, password, databaseName)
}
