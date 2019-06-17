package services

import (
	"github.com/google/uuid"
	"github.com/splisson/devopstic/entities"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

var (
	testEvent = entities.Event{
		Timestamp:   time.Now(),
		PipelineId:  "test",
		Status:      "success",
		CommitId:    uuid.New().String(),
		Environment: "dev",
	}
)

func TestEventSequence(t *testing.T) {

	eventService := NewEventService(testEventStore)

	rand.Seed(time.Now().UnixNano())
	random := rand.Intn(10)
	mult := time.Duration(-5 * random)
	t.Run("should create event commit type", func(t *testing.T) {
		newEvent := testEvent
		newEvent.Timestamp = time.Now().Add(mult * time.Minute)
		newEvent.Type = entities.EVENT_COMMIT
		event, err := eventService.CreateEvent(newEvent)
		assert.Nil(t, err, "no error")
		assert.Equal(t, entities.EVENT_COMMIT, event.Type)
	})
}
