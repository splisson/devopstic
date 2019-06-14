package services

import (
	"github.com/google/uuid"
	"github.com/splisson/devopstic/entities"
	"github.com/stretchr/testify/assert"
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

	commitService := NewCommitService(testCommitStore, nil)

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
		commitEvent := entities.CommitEvent{
			Status:     entities.STATUS_SUCCESS,
			Type:       entities.COMMIT_EVENT_SUBMIT,
			CommitId:   newCommit.CommitId,
			PipelineId: newCommit.PipelineId,
		}
		commit, err := commitService.UpdateCommitByEvent(commitEvent)
		assert.Nil(t, err, "no error")
		assert.Equal(t, entities.COMMIT_STATE_SUBMITTED, commit.State)
	})

	t.Run("should update state from committed to approved", func(t *testing.T) {
		newCommit := testCommit
		newCommit.CommitId = uuid.New().String()
		commit, err := commitService.CreateCommit(newCommit)
		assert.Nil(t, err, "no error")
		assert.Equal(t, entities.COMMIT_STATE_COMMITTED, commit.State)
		commitEvent := entities.CommitEvent{
			Status:     entities.STATUS_SUCCESS,
			Type:       entities.COMMIT_EVENT_APPROVE,
			CommitId:   newCommit.CommitId,
			PipelineId: newCommit.PipelineId,
		}
		commit, err = commitService.UpdateCommitByEvent(commitEvent)
		assert.Nil(t, err, "no error")
		assert.Equal(t, entities.COMMIT_STATE_APPROVED, commit.State)

	})

	t.Run("should update state from submitted to approved", func(t *testing.T) {
		newCommit := testCommit
		newCommit.CommitId = uuid.New().String()
		commit, err := commitService.CreateCommit(newCommit)
		assert.Nil(t, err, "no error")
		assert.Equal(t, entities.COMMIT_STATE_COMMITTED, commit.State)
		commitEvent := entities.CommitEvent{
			Status:     entities.STATUS_SUCCESS,
			Type:       entities.COMMIT_EVENT_SUBMIT,
			CommitId:   newCommit.CommitId,
			PipelineId: newCommit.PipelineId,
		}
		commit, err = commitService.UpdateCommitByEvent(commitEvent)
		assert.Nil(t, err, "no error")
		assert.Equal(t, entities.COMMIT_STATE_SUBMITTED, commit.State)
		commitEvent = entities.CommitEvent{
			Status:     entities.STATUS_SUCCESS,
			Type:       entities.COMMIT_EVENT_APPROVE,
			CommitId:   newCommit.CommitId,
			PipelineId: newCommit.PipelineId,
		}
		commit, err = commitService.UpdateCommitByEvent(commitEvent)
		assert.Nil(t, err, "no error")
		assert.Equal(t, entities.COMMIT_STATE_APPROVED, commit.State)

	})

	t.Run("should update state from committed to deployed", func(t *testing.T) {
		newCommit := testCommit
		newCommit.CommitId = uuid.New().String()
		commit, err := commitService.CreateCommit(newCommit)
		assert.Nil(t, err, "no error")
		assert.Equal(t, entities.COMMIT_STATE_COMMITTED, commit.State)
		commitEvent := entities.CommitEvent{
			Status:     entities.STATUS_SUCCESS,
			Type:       entities.COMMIT_EVENT_DEPLOY,
			CommitId:   newCommit.CommitId,
			PipelineId: newCommit.PipelineId,
		}
		commit, err = commitService.UpdateCommitByEvent(commitEvent)
		assert.Nil(t, err, "no error")
		assert.Equal(t, entities.COMMIT_STATE_DEPLOYED, commit.State)
		deployments, err := testDeploymentStore.GetDeploymentsByCommitId(commit.PipelineId, commit.CommitId)
		deployment := deployments[0]
		assert.Nil(t, err, "no error")
		assert.NotNil(t, deployment, "deployment found")
		assert.Equal(t, commitEvent.Status, deployment.Status)
		assert.Equal(t, commitEvent.CommitId, deployment.CommitId)
	})
}
