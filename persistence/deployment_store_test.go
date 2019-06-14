package persistence

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateGetEvent(t *testing.T) {

	t.Run("should create and get event from db", func(t *testing.T) {
		newDeployment := testDeployment
		deployment, err := testDeploymentStore.CreateDeployment(newDeployment)
		assert.Nil(t, err, "no error")
		assert.NotNil(t, deployment.ID, "id should not be nil")
		assert.NotEmpty(t, deployment.ID, "id should not be empty")
		deployments, err := testDeploymentStore.GetAllDeployments()
		assert.Nil(t, err, "no error")
		assert.NotNil(t, deployments, "deployments exists")
		assert.NotEmpty(t, deployments, "list not empty")
	})
}
