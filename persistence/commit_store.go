package persistence

import (
	"github.com/jinzhu/gorm"
	"github.com/pborman/uuid"
	"github.com/splisson/devopstic/entities"
)

type CommitStoreInterface interface {
	GetCommits() ([]entities.Commit, error)
	GetCommitByPipelineIdAndCommitId(pipelineId string, commitId string) (*entities.Commit, error)
	CreateCommit(event entities.Commit) (*entities.Commit, error)
	UpdateCommit(event entities.Commit) (*entities.Commit, error)
}

type CommitStoreDB struct {
	db *gorm.DB
}

func NewCommitStoreDB(db *gorm.DB) *CommitStoreDB {
	store := new(CommitStoreDB)
	store.db = db
	db.LogMode(true)
	return store
}

func (s *CommitStoreDB) GetCommits() ([]entities.Commit, error) {
	commits := []entities.Commit{}
	db := s.db.Table("commits").Select("*")
	db = db.Find(&commits)
	return commits, db.Error
}

func (s *CommitStoreDB) GetCommitByPipelineIdAndCommitId(pipelineId string, commitId string) (*entities.Commit, error) {
	commit := entities.Commit{}
	db := s.db.Table("commits").Select("*").Where("pipeline_id  = ? AND commit_id = ?", pipelineId, commitId)
	db = db.Find(&commit)
	return &commit, db.Error
}

func (s *CommitStoreDB) CreateCommit(commit entities.Commit) (*entities.Commit, error) {
	commit.ID = uuid.New()

	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return nil, err
	}

	if err := tx.Create(&commit).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return &commit, tx.Commit().Error
}

func (s *CommitStoreDB) UpdateCommit(commit entities.Commit) (*entities.Commit, error) {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return nil, err
	}

	if err := tx.Save(&commit).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return &commit, tx.Commit().Error
}
