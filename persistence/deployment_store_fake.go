package persistence

import (
	"github.com/google/uuid"
	"github.com/splisson/devopstic/entities"
	"time"
)

var (
	testDeployment = entities.Deployment{
		Timestamp:   time.Now(),
		PipelineId:  "test_pipeline",
		Status:      "success",
		CommitId:    "1234567890",
		Environment: "unit_test",
	}
)

type DeploymentStoreFake struct {
	deployments []entities.Deployment
}

func NewDeploymentStoreFake() *DeploymentStoreFake {
	store := new(DeploymentStoreFake)
	store.deployments = make([]entities.Deployment, 0)
	store.CreateDeployment(testDeployment)
	store.CreateDeployment(testDeployment)
	return store
}

func (s *DeploymentStoreFake) GetDeploymentsByCommitId(pipelineId string, commitId string) ([]entities.Deployment, error) {
	deployments := make([]entities.Deployment, 0)
	for _, deployment := range s.deployments {
		if deployment.CommitId == commitId && deployment.PipelineId == pipelineId {
			deployments = append(deployments, deployment)
		}
	}
	return deployments, nil
}

func (s *DeploymentStoreFake) GetAllDeployments() ([]entities.Deployment, error) {
	return s.deployments, nil
}

func (s *DeploymentStoreFake) CreateDeployment(deployment entities.Deployment) (*entities.Deployment, error) {
	deployment.ID = uuid.New().String()
	s.deployments = append(s.deployments, deployment)
	return &deployment, nil
}
