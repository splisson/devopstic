package persistence

import (
	"github.com/splisson/opstic/entities"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)
var (
	testEventStore *EventStoreDB
	testEvent = entities.Event{
		Category: entities.EVENT_CATEGORY_DEPLOY,
		Timestamp: time.Now(),
		PipelineId:		"test_pipeline",
		Status:			"success",
		Commit:			"1234567890",
		Environment:     "unit_test",
	}

)

func TestCreateGetEvent(t *testing.T) {

	t.Run("should create and get event from db", func(t *testing.T) {
		newEvent := testEvent
		event, err := testEventStore.CreateEvent(newEvent)
		assert.Nil(t, err, "no error")
		assert.NotNil(t, event.ID, "id should not be nil")
		assert.NotEmpty(t, event.ID, "id should not be empty")
		events, err := testEventStore.GetEvents()
		assert.Nil(t, err, "no error")
		assert.NotNil(t, events, "events exists")
		assert.NotEmpty(t, events, "list not empty")
	})
}