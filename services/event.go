package services

import (
	"github.com/prometheus/common/log"
	"github.com/splisson/opstic/entities"
	"github.com/splisson/opstic/persistence"
)

// TODO: when creating a deploy successful event, one should record Lead time for change
//  for this commit by finding the build time event for this commit if any received
//event.LeadTime

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
	return s.eventStore.GetEvents()
}

func (s *EventService) CreateEvent(event entities.Event) (*entities.Event, error) {
	if event.Category == entities.EVENT_CATEGORY_DEPLOY &&
		event.Status == entities.EVENT_STATUS_SUCCESS {
		buildEvent, err := s.eventStore.GetEventByCommitAndCategory(event.Commit, entities.EVENT_CATEGORY_BUILD)
		if err != nil {
			log.Infof("no build event for that commit %v", err)
		} else {
			event.LeadTime = event.Timestamp.Unix() - buildEvent.Timestamp.Unix()
		}
	}
	return s.eventStore.CreateEvent(event)
}
