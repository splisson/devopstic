package persistence

import (
	"github.com/google/uuid"
	"github.com/splisson/devopstic/entities"
	"time"
)

var (
	testEvent = entities.Event{
		Timestamp:   time.Now(),
		PipelineId:  "test_pipeline",
		Status:      "success",
		CommitId:    "1234567890",
		Environment: "unit_test",
	}
)

type EventStoreFake struct {
	events []entities.Event
}

func NewEventStoreFake() *EventStoreFake {
	store := new(EventStoreFake)
	store.events = make([]entities.Event, 0)
	store.CreateEvent(testEvent)
	store.CreateEvent(testEvent)
	return store
}

func (s *EventStoreFake) GetEventsByCommitId(pipelineId string, commitId string) ([]entities.Event, error) {
	events := make([]entities.Event, 0)
	for _, event := range s.events {
		if event.CommitId == commitId && event.PipelineId == pipelineId {
			events = append(events, event)
		}
	}
	return events, nil
}

func (s *EventStoreFake) GetAllEvents() ([]entities.Event, error) {
	return s.events, nil
}

func (s *EventStoreFake) CreateEvent(event entities.Event) (*entities.Event, error) {
	event.ID = uuid.New().String()
	s.events = append(s.events, event)
	return &event, nil
}
