package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/splisson/opstic/entities"
	"github.com/splisson/opstic/persistence"
	"github.com/splisson/opstic/representations"
	"github.com/splisson/opstic/services"
)

func GetEvents(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "events",
	})
}

func eventRepresentationToEvent(eventRepresentation representations.Event) entities.Event {
	event := entities.Event{
		Category: eventRepresentation.Category,
		Timestamp : eventRepresentation.Timestamp,
		PipelineId: eventRepresentation.PipelineId,
		Status: eventRepresentation.Status,
		Commit: eventRepresentation.Commit,
		Environment: eventRepresentation.Environment,
	}
	return event
}

func PostEvents(c *gin.Context) {
	db := persistence.NewPostgresqlConnectionWithEnv()
	eventStore := persistence.NewEventDBStore(db)
	eventService := services.NewEventService(eventStore)
	var bodyBytes []byte
	var newEvent entities.Event
	c.Request.Body.Read(bodyBytes)
	var newEventVals representations.Event
	if err := c.Bind(&newEventVals); err != nil {
		c.JSON(400, gin.H{ "error": err })
	}
	//err := json.Unmarshal(bodyBytes, newEvent)
	newEvent = eventRepresentationToEvent(newEventVals)
	event, err := eventService.CreateEvent(newEvent)
	if err != nil {
		c.JSON(500, gin.H{ "error": err })
	}
	c.JSON(200, event)

	//gin.H{
	//		"message": "created",
	//	}
}