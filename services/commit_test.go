package services

import (
	"github.com/google/uuid"
	"github.com/splisson/devopstic/entities"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

var (
	testCommit = entities.Commit{
		CommitTime: time.Now(),
		PipelineId: uuid.New().String(),
		CommitId:   "1234567890",
	}
)

func TestCreateCommit(t *testing.T) {

	commitService := NewCommitService(testCommitStore)

	t.Run("should create build without same pipelineId and commit", func(t *testing.T) {
		newCommit := testCommit
		newCommit.CommitId = "123success"
		newCommit.ApprovalTime = time.Now()
		commit, err := commitService.CreateCommit(newCommit)
		assert.Nil(t, err, "no error")
		newCommit = testCommit
		newCommit.CommitId = "othercommit"
		newCommit.ApprovalTime = time.Now()
		commit, err = commitService.CreateCommit(newCommit)
		assert.Nil(t, err, "no error")
		assert.NotNil(t, commit.ID, "id should not be nil")
		assert.NotEmpty(t, commit.ID, "id should not be empty")
		newCommit = testCommit
		newCommit.CommitId = "123success"
		newCommit.ApprovalTime = time.Now()
		commit, err = commitService.CreateCommit(newCommit)
		assert.NotNil(t, err, "should error because commit already exist")
		assert.Nil(t, commit, "commit should be nil")
	})
}

func TestCommitUpdate(t *testing.T) {

	commitService := NewCommitService(testCommitStore)
	rand.Seed(time.Now().UnixNano())
	random := rand.Intn(10)
	mult := time.Duration(-5 * random)

	t.Run("should create commit in commited state", func(t *testing.T) {
		newCommit := testCommit
		commit, err := commitService.CreateCommit(newCommit)
		assert.Nil(t, err, "no error")
		assert.Equal(t, entities.COMMIT_STATE_COMMITTED, commit.State)
	})

	t.Run("should not create duplicate commit", func(t *testing.T) {
		newCommit := testCommit
		commit, err := commitService.CreateCommit(newCommit)
		assert.NotNil(t, err, "should error because commit already exist")
		assert.Nil(t, commit, "commit should be nil")
	})

	t.Run("should update state to submitted", func(t *testing.T) {
		newCommit := testCommit
		event := entities.Event{
			Status:     entities.STATUS_SUCCESS,
			Type:       entities.EVENT_SUBMIT,
			CommitId:   newCommit.CommitId,
			PipelineId: newCommit.PipelineId,
			Timestamp:  time.Now(),
		}
		commit, err := commitService.UpdateCommitByEvent(event)
		assert.Nil(t, err, "no error")
		assert.Equal(t, entities.COMMIT_STATE_SUBMITTED, commit.State)
	})

	t.Run("should update state from committed to approved", func(t *testing.T) {
		newCommit := testCommit
		newCommit.CommitId = uuid.New().String()
		newCommit.CommitTime = time.Now().Add(mult * time.Minute)
		commit, err := commitService.CreateCommit(newCommit)
		assert.Nil(t, err, "no error")
		assert.Equal(t, entities.COMMIT_STATE_COMMITTED, commit.State)
		event := entities.Event{
			Status:     entities.STATUS_SUCCESS,
			Type:       entities.EVENT_APPROVE,
			CommitId:   newCommit.CommitId,
			PipelineId: newCommit.PipelineId,
			Timestamp:  time.Now(),
		}
		commit, err = commitService.UpdateCommitByEvent(event)
		assert.Nil(t, err, "no error")
		assert.Equal(t, entities.COMMIT_STATE_APPROVED, commit.State)
		assert.True(t, commit.ReviewLeadTime > 0, "review lead time is > 0")
		assert.True(t, commit.TotalLeadTime == commit.ReviewLeadTime, "total lead time = review lead time")
	})

	t.Run("should update state from submitted to approved", func(t *testing.T) {
		newCommit := testCommit
		newCommit.CommitId = uuid.New().String()
		newCommit.CommitTime = time.Now().Add(mult * time.Minute)
		commit, err := commitService.CreateCommit(newCommit)
		assert.Nil(t, err, "no error")
		assert.Equal(t, entities.COMMIT_STATE_COMMITTED, commit.State)
		event := entities.Event{
			Status:     entities.STATUS_SUCCESS,
			Type:       entities.EVENT_SUBMIT,
			CommitId:   newCommit.CommitId,
			PipelineId: newCommit.PipelineId,
			Timestamp:  time.Now(),
		}
		commit, err = commitService.UpdateCommitByEvent(event)
		assert.Nil(t, err, "no error")
		assert.Equal(t, entities.COMMIT_STATE_SUBMITTED, commit.State)
		event = entities.Event{
			Status:     entities.STATUS_SUCCESS,
			Type:       entities.EVENT_APPROVE,
			CommitId:   newCommit.CommitId,
			PipelineId: newCommit.PipelineId,
			Timestamp:  time.Now(),
		}
		commit, err = commitService.UpdateCommitByEvent(event)
		assert.Nil(t, err, "no error")
		assert.Equal(t, entities.COMMIT_STATE_APPROVED, commit.State)
		assert.True(t, commit.ReviewLeadTime > 0, "review lead time is > 0")
		assert.True(t, commit.TotalLeadTime == commit.ReviewLeadTime, "total lead time = review lead time")
	})

	t.Run("should update state from committed to deployed", func(t *testing.T) {
		newCommit := testCommit
		newCommit.CommitId = uuid.New().String()
		newCommit.CommitTime = time.Now().Add(mult * time.Minute)
		commit, err := commitService.CreateCommit(newCommit)
		assert.Nil(t, err, "no error")
		assert.Equal(t, entities.COMMIT_STATE_COMMITTED, commit.State)
		event := entities.Event{
			Status:     entities.STATUS_SUCCESS,
			Type:       entities.EVENT_DEPLOY,
			CommitId:   newCommit.CommitId,
			PipelineId: newCommit.PipelineId,
			Timestamp:  time.Now(),
		}
		commit, err = commitService.UpdateCommitByEvent(event)
		assert.Nil(t, err, "no error")
		assert.Equal(t, entities.COMMIT_STATE_DEPLOYED, commit.State)
		assert.True(t, commit.DeploymentLeadTime > 0, "deployment lead time is > 0")
		assert.True(t, commit.TotalLeadTime == (commit.ReviewLeadTime+commit.DeploymentLeadTime), "total lead time = review lead time + deployment lead time ")
	})
}
