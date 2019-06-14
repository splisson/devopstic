package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/splisson/devopstic/entities"
	"github.com/splisson/devopstic/representations"
	"github.com/splisson/devopstic/services"
)

type DeploymentHandlers struct {
	deploymentService services.DeploymentServiceInterface
}

func NewDeploymentHandlers(eventService services.DeploymentServiceInterface) *DeploymentHandlers {
	handler := new(DeploymentHandlers)
	handler.deploymentService = eventService
	return handler
}

//func representationToDeployment(representation representations.Deployment) entities.Deployment {
//	//leadTime, _ := strconv.ParseInt(representation.LeadTime, 10, 64)
//	timestamp := time.Unix(representation.Timestamp, 0)
//	deployment := entities.Deployment{
//		Model:       entities.Model{ ID: representation.Id },
//		Timestamp:   timestamp,
//		PipelineId:  representation.PipelineId,
//		Status:      representation.Status,
//		CommitId:    representation.CommitId,
//		Environment: representation.Environment,
//	}
//	return deployment
//}

func deploymentToRepresentation(deployment entities.Deployment) representations.Deployment {
	deploymentRepresentation := representations.Deployment{
		Id:          deployment.ID,
		Timestamp:   deployment.Timestamp.Unix(),
		PipelineId:  deployment.PipelineId,
		Status:      deployment.Status,
		CommitId:    deployment.CommitId,
		Environment: deployment.Environment,
	}
	return deploymentRepresentation
}

func (e *DeploymentHandlers) GetDeployments(c *gin.Context) {
	events, err := e.deploymentService.GetDeployments()
	if err != nil {
		c.JSON(500, gin.H{"error": err})
	}
	deploymentList := make([]representations.Deployment, 0)
	for _, item := range events {
		deploymentList = append(deploymentList, deploymentToRepresentation(item))
	}
	results := representations.DeploymentResults{
		Items: deploymentList,
		Count: len(deploymentList),
		Skip:  0,
		Limit: -1,
	}
	c.JSON(200, results)
}
