package persistence

import (
	"errors"
	"github.com/google/uuid"
	"github.com/splisson/devopstic/entities"
	"time"
)

var (
	testCommit = entities.Commit{
		ApprovalTime: time.Now(),
		PipelineId:   "test_pipeline_" + uuid.New().String(),
		CommitId:     uuid.New().String(),
	}
)

type CommitStoreFake struct {
	commits []entities.Commit
}

func NewCommitStoreFake() *CommitStoreFake {
	store := new(CommitStoreFake)
	store.commits = make([]entities.Commit, 0)
	store.CreateCommit(testCommit)
	testCommit.CommitId = uuid.New().String()
	store.CreateCommit(testCommit)
	return store
}

func (s *CommitStoreFake) GetCommits() ([]entities.Commit, error) {
	return s.commits, nil
}

func (s *CommitStoreFake) GetCommitByPipelineIdAndCommitId(pipelineId string, commitId string) (*entities.Commit, error) {
	for _, commit := range s.commits {
		if commit.CommitId == commitId && commit.PipelineId == pipelineId {
			return &commit, nil
		}
	}
	return nil, errors.New("no commit found that matches criteria")
}

func (s *CommitStoreFake) CreateCommit(commit entities.Commit) (*entities.Commit, error) {
	commit.ID = uuid.New().String()
	s.commits = append(s.commits, commit)
	return &commit, nil
}

func (s *CommitStoreFake) UpdateCommit(commit entities.Commit) (*entities.Commit, error) {
	for index, item := range s.commits {
		if item.CommitId == commit.CommitId && item.PipelineId == commit.PipelineId {
			s.commits[index] = commit
			return &commit, nil
		}
	}
	return nil, errors.New("commit not found")
}
