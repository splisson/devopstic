package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/splisson/devopstic/entities"
	"github.com/splisson/devopstic/representations"
	"github.com/splisson/devopstic/services"
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
