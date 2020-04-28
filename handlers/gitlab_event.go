package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/common/log"
	"github.com/splisson/devopstic/entities"
	"github.com/splisson/devopstic/representations"
	"github.com/splisson/devopstic/services"
	"time"
)

type GitlabEventHandlers struct {
	eventService    services.EventServiceInterface
	commitService   services.CommitServiceInterface
	incidentService services.IncidentServiceInterface
}

func NewGitlabEventHandlers(eventService services.EventServiceInterface, commitService services.CommitServiceInterface, incidentService services.IncidentServiceInterface) *GitlabEventHandlers {
	handler := new(GitlabEventHandlers)
	handler.eventService = eventService
	handler.commitService = commitService
	handler.incidentService = incidentService
	return handler
}

func representationGitlabPushEventToEvent(representation representations.GitlabPushEvent) entities.Event {
	timestamp := time.Now()
	log.Infof("received gitlab event %s %s $s", representation.After, representation.Repository.Name, representation.Repository.Name)
	event := entities.Event{
		PipelineId:    representation.Repository.URL,
		CommitId:      representation.After,
		PullRequestId: 0,
		IncidentId:    "",
		Environment:   "",
		Status:        "success",
		Type:          entities.EVENT_COMMIT,
		Timestamp:     timestamp,
	}

	return event
}

func representationGitlabPullRequestEventToEvent(representation representations.GitlabPullRequestEvent) entities.Event {
	timestamp := time.Now()
	log.Infof("received gitlab event %s %s", representation.PullRequest.LastCommit.Id, representation.Repository.Name)
	event := entities.Event{
		PipelineId:    representation.Repository.Name,
		PullRequestId: representation.PullRequest.Id,
		IncidentId:    "",
		Environment:   "",
		Status:        "success",
		Type:          entities.EVENT_SUBMIT,
		Timestamp:     timestamp,
	}
	if representation.PullRequest.Action == "open" {
		event.Type = entities.EVENT_SUBMIT
		event.CommitId = representation.PullRequest.LastCommit.Id
	} else if representation.PullRequest.Action == "close" || representation.PullRequest.Action == "merge" {
		if representation.PullRequest.Action == "close" {
			representation.PullRequest.Merged = true
		}
		event.CommitId = representation.PullRequest.LastCommit.Id
		if representation.PullRequest.Merged {
			// Closing and merging PR = Approve
			event.Type = entities.EVENT_APPROVE
		} else {
			// Closing PR but not merging it
			event.Type = entities.EVENT_SUBMIT
		}
	}
	return event
}

func decodeGitlabEvent(c *gin.Context) (*entities.Event, error) {
	var err error
	newEvent := entities.Event{}

	// Extract name of event
	name := c.GetHeader("X-Gitlab-Event")

	switch name {
	//case "ping":
	//	newEvent = entities.Event{
	//		Type: "ping",
	//	}
	//	break
	case "Push Hook":
		eventValues := representations.GitlabPushEvent{}
		err = c.Bind(&eventValues)
		if err != nil {
			return nil, err
		}
		newEvent = representationGitlabPushEventToEvent(eventValues)
		break

	case "Merge Request Hook":
		eventValues := representations.GitlabPullRequestEvent{}
		err = c.Bind(&eventValues)
		if err != nil {
			return nil, err
		}
		newEvent = representationGitlabPullRequestEventToEvent(eventValues)
		break
	default:
		return nil, errors.New("unsupported event")

	}

	return &newEvent, err
}

func (e *GitlabEventHandlers) PostGitlabEvents(c *gin.Context) {

	newEvent, bindErr := decodeGitlabEvent(c)
	if bindErr != nil {
		c.JSON(400, gin.H{"error": bindErr})
		return
	}

	if newEvent.Type == "ping" {
		c.JSON(200, "ok")
		return
	}

	event, err := e.eventService.CreateEvent(*newEvent)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if event.Type == entities.EVENT_INCIDENT_STATUS_CHANGE {
		_, err = e.incidentService.HandleEvent(*event)
	} else {
		_, _, err = e.commitService.HandleEvent(*event)
	}
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, "ok") //eventToRepresentation(*event))
}
