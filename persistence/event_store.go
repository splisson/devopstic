package persistence

import (
	"github.com/jinzhu/gorm"
	"github.com/pborman/uuid"
	"github.com/splisson/devopstic/entities"
)

type EventStoreInterface interface {
	GetAllEvents() ([]entities.Event, error)
	GetEventsByCommitId(pipelineId string, commitId string) ([]entities.Event, error)
	CreateEvent(event entities.Event) (*entities.Event, error)
}

type EventStoreDB struct {
	db *gorm.DB
}

func NewEventStoreDB(db *gorm.DB) *EventStoreDB {
	store := new(EventStoreDB)
	store.db = db
	db.LogMode(true)
	return store
}

func (s *EventStoreDB) GetAllEvents() ([]entities.Event, error) {
	events := []entities.Event{}
	db := s.db.Table("events").Select("*")
	db = db.Find(&events)
	return events, db.Error
}

func (s *EventStoreDB) GetEventsByCommitId(pipelineId string, commitId string) ([]entities.Event, error) {
	events := []entities.Event{}
	db := s.db.Table("events").Select("*").Where("pipeline_id = ? AND commit = ?", pipelineId, commitId)
	db = db.Find(&events)
	return events, db.Error
}

func (s *EventStoreDB) CreateEvent(event entities.Event) (*entities.Event, error) {
	event.ID = uuid.New()

	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return nil, err
	}

	if err := tx.Create(&event).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return &event, tx.Commit().Error
}
