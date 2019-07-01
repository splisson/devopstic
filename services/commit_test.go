package services

import (
	"github.com/pborman/uuid"
	"github.com/splisson/devopstic/entities"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

var (
	testCommit = entities.Commit{
		CommitTime: time.Now(),
		PipelineId: "test",
		CommitId:   "1234567890",
	}
)

func TestCreateCommit(t *testing.T) {

	commitService := NewCommitService(testCommitStore, testDeploymentStore)

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

	commitService := NewCommitService(testCommitStore, testDeploymentStore)
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
		newCommit.CommitId = uuid.New()
		commit, err := commitService.CreateCommit(newCommit)
		event := entities.Event{
			Status:     entities.STATUS_SUCCESS,
			Type:       entities.EVENT_SUBMIT,
			CommitId:   commit.CommitId,
			PipelineId: commit.PipelineId,
			Timestamp:  time.Now(),
		}
		commit, _, err = commitService.UpdateCommitByEvent(*commit, event)
		assert.Nil(t, err, "no error")
		assert.Equal(t, entities.COMMIT_STATE_SUBMITTED, commit.State)
	})

	t.Run("should update state from committed to approved", func(t *testing.T) {
		newCommit := testCommit
		newCommit.CommitId = uuid.New()
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
		commit, _, err = commitService.HandleEvent(event)
		assert.Nil(t, err, "no error")
		assert.Equal(t, entities.COMMIT_STATE_APPROVED, commit.State)
		assert.True(t, commit.ReviewLeadTime > 0, "review lead time is > 0")
	})

	t.Run("should update state from submitted to approved", func(t *testing.T) {
		newCommit := testCommit
		newCommit.CommitId = uuid.New()
		newCommit.CommitTime = time.Now().Add(2 * mult * time.Minute)
		commit, err := commitService.CreateCommit(newCommit)
		assert.Nil(t, err, "no error")
		assert.Equal(t, entities.COMMIT_STATE_COMMITTED, commit.State)
		event := entities.Event{
			Status:     entities.STATUS_SUCCESS,
			Type:       entities.EVENT_SUBMIT,
			CommitId:   newCommit.CommitId,
			PipelineId: newCommit.PipelineId,
			Timestamp:  time.Now().Add(mult * time.Minute),
		}
		commit, _, err = commitService.HandleEvent(event)
		assert.Nil(t, err, "no error")
		assert.Equal(t, entities.COMMIT_STATE_SUBMITTED, commit.State)
		event = entities.Event{
			Status:     entities.STATUS_SUCCESS,
			Type:       entities.EVENT_APPROVE,
			CommitId:   newCommit.CommitId,
			PipelineId: newCommit.PipelineId,
			Timestamp:  time.Now(),
		}
		commit, _, err = commitService.HandleEvent(event)
		assert.Nil(t, err, "no error")
		assert.Equal(t, entities.COMMIT_STATE_APPROVED, commit.State)
		assert.True(t, commit.ReviewLeadTime > 0, "review lead time is > 0")
	})

	t.Run("should update state from committed to deployed", func(t *testing.T) {
		newCommit := testCommit
		newCommit.CommitId = uuid.New()
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
		commit, _, err = commitService.HandleEvent(event)
		assert.Nil(t, err, "no error")
		assert.Equal(t, entities.COMMIT_STATE_DEPLOYED, commit.State)
	})

	t.Run("should support multiple deployments", func(t *testing.T) {
		var deployment *entities.Deployment
		newCommit := testCommit
		newCommit.CommitId = uuid.New()
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
		event.Environment = "prod"
		commit, deployment, err = commitService.HandleEvent(event)
		assert.Nil(t, err, "no error")
		assert.Equal(t, entities.COMMIT_STATE_DEPLOYED, commit.State)
		previousLeadTime := deployment.LeadTime
		event.Environment = "prod"
		event.Timestamp = time.Now().Add(2 * time.Minute)
		commit, deployment, err = commitService.HandleEvent(event)
		assert.Nil(t, err, "no error")
		assert.Equal(t, entities.COMMIT_STATE_DEPLOYED, commit.State)
		assert.True(t, deployment.LeadTime > previousLeadTime, "deployment lead time increased")

	})
}
