package persistence

import (
	"github.com/docker/distribution/uuid"
	"github.com/jinzhu/gorm"
	"os"
	"testing"
)

func initTestDB(m *testing.M) (*gorm.DB, string) {
	dbId := uuid.Generate().String()
	db, dbFilepath := NewSQLiteConnection(dbId)
	return db, dbFilepath
}

func cleanupTestDB(dbFilepath string) {
	os.Remove(dbFilepath)
}

func TestMain(m *testing.M) {
	db, dbFilepath = initTestDB(m)
	testUserStore = NewUserDBStore(db)
	testEventStore = NewEventDBStore(db)
	CreateTables(db)
	m.Run()
	cleanupTestDB(dbFilepath)
}
