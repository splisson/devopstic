package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/common/log"
	"github.com/splisson/devopstic/entities"
	"github.com/splisson/devopstic/representations"
	"github.com/splisson/devopstic/services"
	"time"
)

type GithubEventHandlers struct {
	eventService    services.EventServiceInterface
	commitService   services.CommitServiceInterface
	incidentService services.IncidentServiceInterface
}

func NewGithubEventHandlers(eventService services.EventServiceInterface, commitService services.CommitServiceInterface, incidentService services.IncidentServiceInterface) *GithubEventHandlers {
	handler := new(GithubEventHandlers)
	handler.eventService = eventService
	handler.commitService = commitService
	handler.incidentService = incidentService
	return handler
}

func representationGithubEventToEvent(representation representations.GithubEvent) entities.Event {
	timestamp := time.Now()
	log.Infof("received github event %s %s $s", representation.Name, representation.Head, representation.Repository.Name)
	event := entities.Event{
		PipelineId:  representation.Repository.Name,
		CommitId:    representation.Head,
		IncidentId:  "",
		Environment: "",
		Status:      "success",
		Type:        entities.EVENT_COMMIT,
		Timestamp:   timestamp,
	}
	if representation.Name == "push" {
		event.Type = entities.EVENT_COMMIT
	}

	return event
}

func (e *GithubEventHandlers) PostGithubEvents(c *gin.Context) {
	var eventValues representations.GithubEvent
	var err error
	if bindErr := c.Bind(&eventValues); bindErr != nil {
		c.JSON(400, gin.H{"error": bindErr})
		return
	}
	newEvent := representationGithubEventToEvent(eventValues)
	event, err := e.eventService.CreateEvent(newEvent)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if event.Type == entities.EVENT_INCIDENT_STATUS_CHANGE {
		_, err = e.incidentService.HandleEvent(*event)
	} else {
		_, err = e.commitService.HandleEvent(*event)
	}
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, "ok") //eventToRepresentation(*event))
}
