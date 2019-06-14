package persistence

import (
	"github.com/jinzhu/gorm"
	"github.com/pborman/uuid"
	"github.com/splisson/devopstic/entities"
)

type DeploymentStoreInterface interface {
	GetAllDeployments() ([]entities.Deployment, error)
	GetDeploymentsByCommitId(pipelineId string, commitId string) ([]entities.Deployment, error)
	CreateDeployment(event entities.Deployment) (*entities.Deployment, error)
}

type DeploymentStoreDB struct {
	db *gorm.DB
}

func NewDeploymentDBStore(db *gorm.DB) *DeploymentStoreDB {
	store := new(DeploymentStoreDB)
	store.db = db
	db.LogMode(true)
	return store
}

func (s *DeploymentStoreDB) GetAllDeployments() ([]entities.Deployment, error) {
	deployments := []entities.Deployment{}
	db := s.db.Table("deployments").Select("*")
	db = db.Find(&deployments)
	return deployments, db.Error
}

func (s *DeploymentStoreDB) GetDeploymentsByCommitId(pipelineId string, commitId string) ([]entities.Deployment, error) {
	deployments := []entities.Deployment{}
	db := s.db.Table("deployments").Select("*").Where("pipeline_id = ? AND commit = ?", pipelineId, commitId)
	db = db.Find(&deployments)
	return deployments, db.Error
}

func (s *DeploymentStoreDB) CreateDeployment(deployment entities.Deployment) (*entities.Deployment, error) {
	deployment.ID = uuid.New()

	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return nil, err
	}

	if err := tx.Create(&deployment).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return &deployment, tx.Commit().Error
}
