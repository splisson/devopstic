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

func representationGithubPushEventToEvent(representation representations.GithubPushEvent) entities.Event {
	timestamp := time.Now()
	log.Infof("received github event %s %s $s", representation.After, representation.Repository.Name, representation.Repository.Name)
	event := entities.Event{
		PipelineId:    representation.Repository.Name,
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

func representationGithubPullRequestEventToEvent(representation representations.GithubPullRequestEvent) entities.Event {
	timestamp := time.Now()
	log.Infof("received github event %s %s $s", representation.PullRequest.MergeCommitSha, representation.Repository.Name, representation.Repository.Name)
	event := entities.Event{
		PipelineId:    representation.Repository.Name,
		PullRequestId: representation.PullRequest.Id,
		IncidentId:    "",
		Environment:   "",
		Status:        "success",
		Type:          entities.EVENT_SUBMIT,
		Timestamp:     timestamp,
	}
	if representation.Action == "opened" {
		event.Type = entities.EVENT_SUBMIT
		event.CommitId = representation.PullRequest.Head.Sha
	} else if representation.Action == "closed" {
		event.CommitId = representation.PullRequest.MergeCommitSha
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

func decodeGithubEvent(c *gin.Context) (*entities.Event, error) {
	var err error
	newEvent := entities.Event{}

	// Extract name of event
	name := c.GetHeader("X-Github-Event")

	switch name {
	case "push":
		eventValues := representations.GithubPushEvent{}
		err = c.Bind(&eventValues)
		if err != nil {
			return nil, err
		}
		newEvent = representationGithubPushEventToEvent(eventValues)
		break

	case "pull_request":
		eventValues := representations.GithubPullRequestEvent{}
		err = c.Bind(&eventValues)
		if err != nil {
			return nil, err
		}
		newEvent = representationGithubPullRequestEventToEvent(eventValues)
		break
	default:
		return nil, errors.New("unsupported event")

	}

	return &newEvent, err
}

func (e *GithubEventHandlers) PostGithubEvents(c *gin.Context) {

	newEvent, bindErr := decodeGithubEvent(c)
	if bindErr != nil {
		c.JSON(400, gin.H{"error": bindErr})
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
		_, err = e.commitService.HandleEvent(*event)
	}
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, "ok") //eventToRepresentation(*event))
}
