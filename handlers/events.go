package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/splisson/opstic/entities"
	"github.com/splisson/opstic/representations"
	"github.com/splisson/opstic/services"
)

type EventHandlers struct {
	eventService services.EventServiceInterface
}

func NewEventHandlers(eventService services.EventServiceInterface) *EventHandlers {
	handler := new(EventHandlers)
	handler.eventService = eventService
	return handler
}

func (e *EventHandlers) GetEvents(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "events",
	})
}

func eventRepresentationToEvent(eventRepresentation representations.Event) entities.Event {
	//leadTime, _ := strconv.ParseInt(eventRepresentation.LeadTime, 10, 64)

	event := entities.Event{
		Category:    eventRepresentation.Category,
		Timestamp:   eventRepresentation.Timestamp,
		PipelineId:  eventRepresentation.PipelineId,
		Status:      eventRepresentation.Status,
		Commit:      eventRepresentation.Commit,
		Environment: eventRepresentation.Environment,
	}
	return event
}

func (e *EventHandlers) PostEvents(c *gin.Context) {
	var newEvent entities.Event

	var newEventVals representations.Event
	if err := c.Bind(&newEventVals); err != nil {
		c.JSON(400, gin.H{"error": err})
	}
	//err := json.Unmarshal(bodyBytes, newEvent)
	newEvent = eventRepresentationToEvent(newEventVals)
	event, err := e.eventService.CreateEvent(newEvent)
	if err != nil {
		c.JSON(500, gin.H{"error": err})
	}
	c.JSON(200, event)

	//gin.H{
	//		"message": "created",
	//	}
}
