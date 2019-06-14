package services

import (
	"github.com/splisson/devopstic/entities"
	"github.com/splisson/devopstic/persistence"
)

type DeploymentServiceInterface interface {
	CreateDeployment(event entities.Deployment) (*entities.Deployment, error)
	GetDeployments() ([]entities.Deployment, error)
}

type DeploymentService struct {
	deploymentStore persistence.DeploymentStoreInterface
	commitStore     persistence.CommitStoreInterface
}

func NewDeploymentService(deploymentStore persistence.DeploymentStoreInterface, commitStore persistence.CommitStoreInterface) *DeploymentService {
	service := new(DeploymentService)
	service.deploymentStore = deploymentStore
	service.commitStore = commitStore
	return service
}

func (s *DeploymentService) GetDeployments() ([]entities.Deployment, error) {
	return s.deploymentStore.GetAllDeployments()
}

func (s *DeploymentService) CreateDeployment(deployment entities.Deployment) (*entities.Deployment, error) {
	return s.deploymentStore.CreateDeployment(deployment)
}
