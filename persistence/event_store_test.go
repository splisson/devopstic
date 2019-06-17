package persistence

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateGetEvent(t *testing.T) {

	t.Run("should create and get event from db", func(t *testing.T) {
		newEvent := testEvent
		event, err := testEventStore.CreateEvent(newEvent)
		assert.Nil(t, err, "no error")
		assert.NotNil(t, event.ID, "id should not be nil")
		assert.NotEmpty(t, event.ID, "id should not be empty")
		events, err := testEventStore.GetAllEvents()
		assert.Nil(t, err, "no error")
		assert.NotNil(t, events, "events exists")
		assert.NotEmpty(t, events, "list not empty")
	})
}
