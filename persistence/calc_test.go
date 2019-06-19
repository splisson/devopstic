package postgres_test

import (
	"github.com/jinzhu/gorm"
	"github.com/splisson/devopstic/persistence"
	"os"
	"testing"
)

var (
	db *gorm.DB
)

func initTestPosgresDB(m *testing.M) *gorm.DB {
	return persistence.CreatePostgresDBConnection("localhost", "5432", os.Getenv("DATABASE_USER"), os.Getenv("DATABASE_PASSWORD"), "opstic")
}

func TestMain(m *testing.M) {
	db = initTestPosgresDB(m)
	persistence.CreateTables(db)
	m.Run()

}

func TestBuilds(t *testing.T) {
	t.Run("should create events", func(t *testing.T) {

	})

}
