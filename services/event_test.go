package services

import (
	"github.com/google/uuid"
	"github.com/splisson/devopstic/entities"
	"time"
)

var (
	testEvent = entities.Event{
		Timestamp:   time.Now(),
		PipelineId:  uuid.New().String(),
		Status:      "success",
		CommitId:    uuid.New().String(),
		Environment: "unit_test",
	}
)
