package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/splisson/devopstic/entities"
	"github.com/splisson/devopstic/representations"
	"github.com/splisson/devopstic/services"
	"time"
)

type CommitHandlers struct {
	commitService services.CommitServiceInterface
}

func NewCommitHandlers(commitService services.CommitServiceInterface) *CommitHandlers {
	handler := new(CommitHandlers)
	handler.commitService = commitService
	return handler
}

func commitToRepresentation(commit entities.Commit) representations.Commit {
	return representations.Commit{
		PipelineId:         commit.PipelineId,
		CommitId:           commit.CommitId,
		Id:                 commit.ID,
		CommitTime:         commit.CommitTime.Unix(),
		SubmitTime:         commit.SubmitTime.Unix(),
		ApprovalTime:       commit.ApprovalTime.Unix(),
		DeploymentTime:     commit.DeploymentTime.Unix(),
		TotalLeadTime:      commit.TotalLeadTime,
		DeploymentLeadTime: commit.DeploymentLeadTime,
		ReviewLeadTime:     commit.ReviewLeadTime,
	}
}
func representationToCommitEvent(representation representations.CommitEvent) entities.CommitEvent {
	timestamp := time.Unix(representation.Timestamp, 0)
	return entities.CommitEvent{
		PipelineId:  representation.PipelineId,
		CommitId:    representation.CommitId,
		Environment: representation.Environment,
		Status:      representation.Status,
		Type:        representation.Type,
		Timestamp:   timestamp,
	}
}

func (e *CommitHandlers) GetCommits(c *gin.Context) {
	events, err := e.commitService.GetCommits()
	if err != nil {
		c.JSON(500, gin.H{"error": err})
	}
	commitList := make([]representations.Commit, 0)
	for _, item := range events {
		commitList = append(commitList, commitToRepresentation(item))
	}
	results := representations.CommitResults{
		Items: commitList,
		Count: len(commitList),
		Skip:  0,
		Limit: -1,
	}
	c.JSON(200, results)
}

func (e *CommitHandlers) PostCommitEvent(c *gin.Context) {
	var commit *entities.Commit
	var commitEventValues representations.CommitEvent
	var err error
	if bindErr := c.Bind(&commitEventValues); bindErr != nil {
		c.JSON(400, gin.H{"error": bindErr})
		return
	}
	commitEvent := representationToCommitEvent(commitEventValues)
	if commitEvent.Type == entities.COMMIT_EVENT_COMMIT {
		// Create
		newCommit := entities.Commit{
			PipelineId: commitEvent.PipelineId,
			CommitId:   commitEvent.CommitId,
			State:      entities.COMMIT_STATE_COMMITTED,
			CommitTime: commitEvent.Timestamp,
		}
		commit, err = e.commitService.CreateCommit(newCommit)

	} else {
		// Update
		commit, err = e.commitService.UpdateCommitByEvent(commitEvent)
	}
	if err != nil {
		c.JSON(500, gin.H{"error": err})
		return
	}
	c.JSON(200, commitToRepresentation(*commit))
}
