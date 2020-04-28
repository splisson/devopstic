package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/common/log"
	"github.com/splisson/devopstic/entities"
	"github.com/splisson/devopstic/representations"
	"github.com/splisson/devopstic/services"
	"io/ioutil"
	"time"
)

type PagerDutyHandlers struct {
	incidentService services.IncidentServiceInterface
	eventService    services.EventServiceInterface
}

func NewPagerDutyHandlers(incidentService services.IncidentServiceInterface, eventService services.EventServiceInterface) *PagerDutyHandlers {
	handler := new(PagerDutyHandlers)
	handler.incidentService = incidentService
	handler.eventService = eventService
	return handler
}

func representationPagerDutyIncidentToEvent(representation representations.PDIncident) entities.Event {
	timestamp := time.Now()
	log.Infof("received pagerduty event %d %s", representation.IncidentNumber, representation.Status)
	event := entities.Event{
		PipelineId:    "",
		CommitId:      "",
		PullRequestId: 0,
		IncidentId:    fmt.Sprintf("%d", representation.IncidentNumber),
		Environment:   "",
		Type:          entities.EVENT_INCIDENT_STATUS_CHANGE,
		Timestamp:     timestamp,
	}

	if representation.Status == representations.PAGERDUTY_INCIDENT_STATUS_TRIGGERED {
		event.Status = entities.STATUS_FAILURE
	} else {
		event.Status = entities.STATUS_SUCCESS
	}

	return event
}

func decodePagerDuty(c *gin.Context) (*entities.Event, error) {
	var err error
	newEvent := entities.Event{}

	eventValues := representations.PagerDutyIncident{}
	err = c.Bind(&eventValues)
	if err != nil {
		return nil, err
	}
	newEvent = representationPagerDutyIncidentToEvent(eventValues.Messages[0].Incident)

	return &newEvent, err
}

func (e *PagerDutyHandlers) PostPagerDutyIncidents(c *gin.Context) {

	bytes, _ := ioutil.ReadAll(c.Request.Body)
	fmt.Printf("%s", bytes)

	//newEvent, bindErr := decodePagerDuty(c)
	//if bindErr != nil {
	//	c.JSON(400, gin.H{"error": bindErr})
	//	return
	//}
	//
	//event, err := e.eventService.CreateEvent(*newEvent)
	//if err != nil {
	//	c.JSON(500, gin.H{"error": err.Error()})
	//	return
	//}
	//_, err = e.incidentService.HandleEvent(*event)
	//if err != nil {
	//	c.JSON(500, gin.H{"error": err.Error()})
	//	return
	//}

	c.JSON(200, gin.H{"status": "success"})
}
