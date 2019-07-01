package persistence

import (
	"github.com/jinzhu/gorm"
	"github.com/pborman/uuid"
	"github.com/splisson/devopstic/entities"
)

type DeploymentStoreInterface interface {
	GetDeploymentByCommitIdAndEnvironment(commitId string, environment string) (*entities.Deployment, error)
	CreateDeployment(deployment entities.Deployment) (*entities.Deployment, error)
	UpdateDeployment(deployment entities.Deployment) (*entities.Deployment, error)
}

type DeploymentStoreDB struct {
	db *gorm.DB
}

func NewDeploymentStoreDB(db *gorm.DB) *DeploymentStoreDB {
	store := new(DeploymentStoreDB)
	store.db = db
	db.LogMode(true)
	return store
}

func (s *DeploymentStoreDB) GetDeploymentByCommitIdAndEnvironment(commitId string, environment string) (*entities.Deployment, error) {
	deployment := entities.Deployment{}
	db := s.db.Table("deployments").Select("*").Where("environment = ? AND commit_id = ?", environment, commitId)
	db = db.Find(&deployment)
	return &deployment, db.Error
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

func (s *DeploymentStoreDB) UpdateDeployment(deployment entities.Deployment) (*entities.Deployment, error) {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return nil, err
	}

	if err := tx.Save(&deployment).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return &deployment, tx.Commit().Error
}
