package services

import (
	"github.com/splisson/devopstic/entities"
	"github.com/splisson/devopstic/persistence"
)

type EventServiceInterface interface {
	CreateEvent(event entities.Event) (*entities.Event, error)
	GetEvents() ([]entities.Event, error)
}

type EventService struct {
	eventStore persistence.EventStoreInterface
}

func NewEventService(eventStore persistence.EventStoreInterface) *EventService {
	service := new(EventService)
	service.eventStore = eventStore
	return service
}

func (s *EventService) GetEvents() ([]entities.Event, error) {
	return s.eventStore.GetAllEvents()
}

func (s *EventService) CreateEvent(event entities.Event) (*entities.Event, error) {
	return s.eventStore.CreateEvent(event)
}
