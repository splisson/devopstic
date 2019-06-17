package persistence

import (
	"github.com/jinzhu/gorm"
	"github.com/pborman/uuid"
	"github.com/splisson/devopstic/entities"
)

type IncidentStoreInterface interface {
	GetIncidents() ([]entities.Incident, error)
	GetIncidentByIncidentId(incidentId string) (*entities.Incident, error)
	GetLatestFailureIncidentByPipelineIdAndEnvironment(pipelineId string, environment string) (*entities.Incident, error)
	CreateIncident(event entities.Incident) (*entities.Incident, error)
	UpdateIncident(event entities.Incident) (*entities.Incident, error)
}

type IncidentStoreDB struct {
	db *gorm.DB
}

func NewIncidentStoreDB(db *gorm.DB) *IncidentStoreDB {
	store := new(IncidentStoreDB)
	store.db = db
	db.LogMode(true)
	return store
}

func (s *IncidentStoreDB) GetIncidents() ([]entities.Incident, error) {
	incidents := []entities.Incident{}
	db := s.db.Table("incidents").Select("*")
	db = db.Find(&incidents)
	return incidents, db.Error
}

func (s *IncidentStoreDB) GetIncidentByIncidentId(incidentId string) (*entities.Incident, error) {
	incident := entities.Incident{}
	db := s.db.Table("incidents").Select("*").Where("incident_id = ?", incidentId)
	db = db.Find(&incident)
	return &incident, db.Error
}

func (s *IncidentStoreDB) GetLatestFailureIncidentByPipelineIdAndEnvironment(pipelineId string, environment string) (*entities.Incident, error) {
	event := entities.Incident{}
	db := s.db.Table("incidents").Select("*").
		Where("status = ? AND pipeline_id = ? AND environment = ?", entities.STATUS_FAILURE, pipelineId, environment).
		Order("timestamp DESC")
	db = db.First(&event)
	return &event, db.Error
}

func (s *IncidentStoreDB) CreateIncident(incident entities.Incident) (*entities.Incident, error) {
	incident.ID = uuid.New()

	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return nil, err
	}

	if err := tx.Create(&incident).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return &incident, tx.Commit().Error
}

func (s *IncidentStoreDB) UpdateIncident(incident entities.Incident) (*entities.Incident, error) {

	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return nil, err
	}

	if err := tx.Save(&incident).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return &incident, tx.Commit().Error
}
