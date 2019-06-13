package persistence

import (
	"errors"
	"github.com/google/uuid"
	"github.com/splisson/opstic/entities"
	"time"
)

var (
	testEvent = entities.Event{
		Category:    entities.EVENT_CATEGORY_DEPLOY,
		Timestamp:   time.Now(),
		PipelineId:  "test_pipeline",
		Status:      "success",
		Commit:      "1234567890",
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

func (s *EventStoreFake) GetEvents() ([]entities.Event, error) {
	//events := []entities.Event{}
	//testEvent.ID = uuid.New().String()
	//events = append(events, testEvent)
	//otherEvent := testEvent
	//otherEvent.ID = uuid.New().String()
	//events = append(events, otherEvent)
	return s.events, nil
}

func (s *EventStoreFake) GetEventByCommitAndCategory(commit string, category string) (*entities.Event, error) {
	for _, event := range s.events {
		if event.Commit == commit && event.Category == category {
			return &event, nil
		}
	}
	return nil, errors.New("no event found that matches criteria")
}

func (s *EventStoreFake) GetLatestFailureEventByPipelineIdAndEnvironment(pipelineId string, environment string) (*entities.Event, error) {
	for _, event := range s.events {
		if event.PipelineId == pipelineId && event.Category == entities.EVENT_CATEGORY_INCIDENT && event.Environment == environment {
			return &event, nil
		}
	}
	return nil, errors.New("no event found that matches criteria")
}

func (s *EventStoreFake) CreateEvent(event entities.Event) (*entities.Event, error) {
	event.ID = uuid.New().String()
	s.events = append(s.events, event)
	return &event, nil
}
