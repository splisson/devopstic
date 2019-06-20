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
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", url, port, username, password, databaseName)
	db, err := gorm.Open("postgres", connectionString)
	if err != nil {
		panic("failed to connect to database")
	}
	return db
}

func NewPostgresqlConnectionLocalhost() *gorm.DB {
	return CreatePostgresDBConnection("localhost", "5432", os.Getenv("DATABASE_USER"), os.Getenv("DATABASE_PASSWORD"), "opstic")
}

func NewPostgresqlConnectionWithEnv() *gorm.DB {
	url := os.Getenv("DATABASE_HOST")
	if url == "" {
		url = "localhost"
	}
	port := os.Getenv("DATABASE_PORT")
	if port == "" {
		port = "5432"
	}
	username := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASSWORD")
	databaseName := os.Getenv("DATABASE_NAME")
	if databaseName == "" {
		databaseName = "opstic"
	}
	return CreatePostgresDBConnection(url, port, username, password, databaseName)
}
