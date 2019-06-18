package services

import (
	"errors"
	"github.com/prometheus/common/log"
	"github.com/splisson/devopstic/entities"
	"github.com/splisson/devopstic/persistence"
	"time"
)

type IncidentServiceInterface interface {
	CreateOrUpdateIncident(incident entities.Incident) (*entities.Incident, error)
	GetIncidents() ([]entities.Incident, error)
	HandleEvent(event entities.Event) (*entities.Incident, error)
}

type IncidentService struct {
	incidentStore persistence.IncidentStoreInterface
}

func NewIncidentService(incidentStore persistence.IncidentStoreInterface) *IncidentService {
	service := new(IncidentService)
	service.incidentStore = incidentStore
	return service
}

func (s *IncidentService) HandleEvent(event entities.Event) (*entities.Incident, error) {
	if event.Type == entities.EVENT_INCIDENT_STATUS_CHANGE {
		newIncident := entities.Incident{
			IncidentId:  event.IncidentId,
			PipelineId:  event.PipelineId,
			Status:      event.Status,
			Environment: event.Environment,
		}
		if event.Status == entities.STATUS_SUCCESS {
			newIncident.ResolutionTime = event.Timestamp
		} else {
			newIncident.OpeningTime = event.Timestamp
		}
		incident, err := s.CreateOrUpdateIncident(newIncident)
		return incident, err
	} else {
		return nil, errors.New("event not supported by incident service")
	}
}

func (s *IncidentService) GetIncidents() ([]entities.Incident, error) {
	return s.incidentStore.GetIncidents()
}

func validateIncident(incident entities.Incident) error {
	if len(incident.IncidentId) == 0 {
		return errors.New("incidentId cannot be empty")
	}
	if len(incident.PipelineId) == 0 {
		return errors.New("pipelineId cannot be empty")
	}
	return nil
}
func (s *IncidentService) CreateOrUpdateIncident(incident entities.Incident) (*entities.Incident, error) {
	if err := validateIncident(incident); err != nil {
		return nil, err
	}
	// Existing incident failure?
	var existingIncident *entities.Incident = nil
	var err error = nil
	if len(incident.IncidentId) > 0 {
		existingIncident, err = s.incidentStore.GetIncidentByIncidentId(incident.IncidentId)
	}
	if err != nil || existingIncident == nil {
		// No existing incident => create
		log.Infof("no incident for that IncidentId => create")
		return s.incidentStore.CreateIncident(incident)
	} else {
		// Existing incident so update based on status
		if existingIncident.Status == entities.STATUS_FAILURE {
			if incident.Status == entities.STATUS_SUCCESS {
				// Recovery
				incident.TimeToRestore = incident.ResolutionTime.Unix() - existingIncident.OpeningTime.Unix()
				incident.ID = existingIncident.ID
				incident.OpeningTime = existingIncident.OpeningTime
				return s.incidentStore.UpdateIncident(incident)
			} else {
				// NOOP: keep original failure time
				return existingIncident, nil
			}
		} else {
			if incident.Status == entities.STATUS_SUCCESS {
				// NOOP: keep first recovery
				return existingIncident, nil
			} else {
				// Updating existing back to failure: reset time to restore
				incident.ResolutionTime = time.Unix(0, 0)
				incident.TimeToRestore = 0
				return s.incidentStore.UpdateIncident(incident)
			}
		}
	}

}
