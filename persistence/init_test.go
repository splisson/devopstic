package persistence

import (
	"github.com/docker/distribution/uuid"
	"github.com/jinzhu/gorm"
	"os"
	"testing"
)

var (
	db         *gorm.DB
	dbFilepath string

	testEventStore    *EventStoreDB
	testCommitStore   *CommitStoreDB
	testIncidentStore *IncidentStoreDB
	testUserStore     *UserDBStore
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
	testEventStore = NewEventStoreDB(db)
	testCommitStore = NewCommitStoreDB(db)
	testIncidentStore = NewIncidentStoreDB(db)
	CreateTables(db)
	exitCode := m.Run()
	cleanupTestDB(dbFilepath)
	os.Exit(exitCode)
}
