package postgres_test

import (
	"github.com/jinzhu/gorm"
	"github.com/prometheus/common/log"
	"github.com/splisson/opstic/entities"
	"github.com/splisson/opstic/persistence"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var (
	db *gorm.DB
)

func initTestPosgresDB(m *testing.M) (*gorm.DB) {
	return persistence.CreatePostgresDBConnection("localhost", "5432", "postgres", "w3yv", "opstic")
}


func TestMain(m *testing.M) {
	db = initTestPosgresDB(m)
	persistence.CreateTables(db)
	m.Run()

}

func TestEvents( t *testing.T){
	t.Run("should create events", func(t *testing.T) {
		store := persistence.NewEventDBStore(db)
		newEvent := entities.Event{
			Timestamp: time.Now(),
			Category: entities.EVENT_CATEGORY_DEPLOY,
			Environment: "dev",
			PipelineId: "api-pipeline",
			Status: "success",
			Commit: "1234567890",
		}
		event, err := store.CreateEvent(newEvent)
		assert.Nil(t, err, "no error")
		assert.NotNil(t, event.ID, "id should not be nil")
		assert.NotEmpty(t, event.ID, "id should not be empty")

	})

}

func TestGetEvents(t *testing.T) {
	store := persistence.NewEventDBStore(db)
	events, err := store.GetEvents()
	if err != nil {
		log.Error(err)
	}
	log.Info(events)
}

