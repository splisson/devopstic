package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/splisson/devopstic/entities"
	"github.com/splisson/devopstic/representations"
	"github.com/splisson/devopstic/services"
	"time"
)

type EventHandlers struct {
	eventService    services.EventServiceInterface
	commitService   services.CommitServiceInterface
	incidentService services.IncidentServiceInterface
}

func NewEventHandlers(eventService services.EventServiceInterface, commitService services.CommitServiceInterface, incidentService services.IncidentServiceInterface) *EventHandlers {
	handler := new(EventHandlers)
	handler.eventService = eventService
	handler.commitService = commitService
	handler.incidentService = incidentService
	return handler
}

func representationToEvent(representation representations.Event) entities.Event {
	timestamp := time.Unix(representation.Timestamp, 0)
	return entities.Event{
		PipelineId:  representation.PipelineId,
		CommitId:    representation.CommitId,
		IncidentId:  representation.IncidentId,
		Environment: representation.Environment,
		Status:      representation.Status,
		Type:        representation.Type,
		Timestamp:   timestamp,
	}
}

func (e *EventHandlers) PostEvents(c *gin.Context) {
	var eventValues representations.Event
	var err error
	if bindErr := c.Bind(&eventValues); bindErr != nil {
		c.JSON(400, gin.H{"error": bindErr})
		return
	}
	newEvent := representationToEvent(eventValues)
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
	c.JSON(200, eventToRepresentation(*event))
}

func eventToRepresentation(event entities.Event) representations.Event {
	eventRepresentation := representations.Event{
		Id:          event.ID,
		Timestamp:   event.Timestamp.Unix(),
		PipelineId:  event.PipelineId,
		Status:      event.Status,
		CommitId:    event.CommitId,
		Environment: event.Environment,
		Type:        event.Type,
	}
	return eventRepresentation
}

func (e *EventHandlers) GetEvents(c *gin.Context) {
	events, err := e.eventService.GetEvents()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	}
	eventList := make([]representations.Event, 0)
	for _, item := range events {
		eventList = append(eventList, eventToRepresentation(item))
	}
	results := representations.EventResults{
		Items: eventList,
		Count: len(eventList),
		Skip:  0,
		Limit: -1,
	}
	c.JSON(200, results)
}
