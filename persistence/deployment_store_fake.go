package persistence

import (
	"errors"
	"github.com/google/uuid"
	"github.com/splisson/devopstic/entities"
	"time"
)

var (
	testDeployment = entities.Deployment{
		Timestamp:   time.Now(),
		CommitId:    "test_pipeline_" + uuid.New().String(),
		Environment: "test",
	}
)

type DeploymentStoreFake struct {
	deployments []entities.Deployment
}

func NewDeploymentStoreFake() *DeploymentStoreFake {
	store := new(DeploymentStoreFake)
	store.deployments = make([]entities.Deployment, 0)
	store.CreateDeployment(testDeployment)
	testDeployment.CommitId = uuid.New().String()
	store.CreateDeployment(testDeployment)
	return store
}

func (s *DeploymentStoreFake) GetDeployments() ([]entities.Deployment, error) {
	return s.deployments, nil
}

func (s *DeploymentStoreFake) GetDeploymentByCommitIdAndEnvironment(commitId string, environment string) (*entities.Deployment, error) {
	for _, deployment := range s.deployments {
		if deployment.CommitId == commitId && deployment.Environment == environment {
			return &deployment, nil
		}
	}
	return nil, errors.New("no deployment found that matches criteria")
}

func (s *DeploymentStoreFake) CreateDeployment(deployment entities.Deployment) (*entities.Deployment, error) {
	deployment.ID = uuid.New().String()
	s.deployments = append(s.deployments, deployment)
	return &deployment, nil
}

func (s *DeploymentStoreFake) UpdateDeployment(deployment entities.Deployment) (*entities.Deployment, error) {
	for index, item := range s.deployments {
		if item.CommitId == deployment.CommitId && item.Environment == deployment.Environment {
			s.deployments[index] = deployment
			return &deployment, nil
		}
	}
	return nil, errors.New("deployment not found")
}
