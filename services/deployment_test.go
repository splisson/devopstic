package services

import (
	"github.com/google/uuid"
	"github.com/splisson/devopstic/entities"
	"time"
)

var (
	testDeployment = entities.Deployment{
		Timestamp:   time.Now(),
		PipelineId:  uuid.New().String(),
		Status:      "success",
		CommitId:    uuid.New().String(),
		Environment: "unit_test",
	}
)

//func TestDeploymentWithLeadTime(t *testing.T) {
//
//	deploymentService := NewDeploymentService(testDeploymentStore, testCommitStore)
//
//	t.Run("success deploy should fill lead time if there is a build for that commit", func(t *testing.T) {
//		newCommit := testCommit
//		newCommit.CommitId = "123success"
//		newCommit.ApprovalTime = time.Now().Add(-5 * time.Minute)
//		commit, err := testCommitStore.CreateCommit(newCommit)
//		assert.Nil(t, err, "no error")
//		assert.NotNil(t, commit, "commit not nil")
//		newDeployment := testDeployment
//		newDeployment.PipelineId = newCommit.PipelineId
//		newDeployment.CommitId = "123success"
//		newDeployment.Timestamp = time.Now()
//		deployment, err := deploymentService.CreateDeployment(newDeployment)
//		assert.Nil(t, err, "no error")
//		assert.NotNil(t, deployment.ID, "id should not be nil")
//		assert.NotEmpty(t, deployment.ID, "id should not be empty")
//		assert.True(t, deployment.LeadTime > 0, "lead time should be > 0")
//	})
//
//
//}
