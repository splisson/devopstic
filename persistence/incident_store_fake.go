package persistence

import (
	"errors"
	"github.com/google/uuid"
	"github.com/splisson/devopstic/entities"
	"time"
)

var (
	testIncident = entities.Incident{
		OpeningTime: time.Now(),
		PipelineId:  "test_pipeline",
		Status:      "success",
		IncidentId:  "1234567890",
		Environment: "unit_test",
	}
)

type IncidentStoreFake struct {
	incidents []entities.Incident
}

func NewIncidentStoreFake() *IncidentStoreFake {
	store := new(IncidentStoreFake)
	store.incidents = make([]entities.Incident, 0)
	store.CreateIncident(testIncident)
	store.CreateIncident(testIncident)
	return store
}

func (s *IncidentStoreFake) GetIncidents() ([]entities.Incident, error) {
	return s.incidents, nil
}

func (s *IncidentStoreFake) GetIncidentByIncidentId(incidentId string) (*entities.Incident, error) {
	for _, incident := range s.incidents {
		if incident.IncidentId == incidentId {
			return &incident, nil
		}
	}
	return nil, errors.New("no incident found that matches criteria")
}

func (s *IncidentStoreFake) GetLatestFailureIncidentByPipelineIdAndEnvironment(pipelineId string, environment string) (*entities.Incident, error) {
	for _, incident := range s.incidents {
		if incident.PipelineId == pipelineId && incident.Status == entities.STATUS_FAILURE && incident.Environment == environment {
			return &incident, nil
		}
	}
	return nil, errors.New("no incident found that matches criteria")
}

func (s *IncidentStoreFake) CreateIncident(incident entities.Incident) (*entities.Incident, error) {
	incident.ID = uuid.New().String()
	s.incidents = append(s.incidents, incident)
	return &incident, nil
}

func (s *IncidentStoreFake) UpdateIncident(incident entities.Incident) (*entities.Incident, error) {
	for index, item := range s.incidents {
		if item.IncidentId == incident.IncidentId {
			if incident.OpeningTime == time.Unix(0, 0) {
				incident.OpeningTime = item.OpeningTime
			}
			if incident.ResolutionTime == time.Unix(0, 0) {
				incident.ResolutionTime = item.ResolutionTime
			}
			s.incidents[index] = incident
			return &incident, nil
		}
	}
	return nil, errors.New("no incident found ")
}
