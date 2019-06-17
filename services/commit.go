package services

import (
	"errors"
	"github.com/splisson/devopstic/entities"
	"github.com/splisson/devopstic/persistence"
)

type CommitServiceInterface interface {
	HandleEvent(event entities.Event) error
	CreateCommit(event entities.Commit) (*entities.Commit, error)
	GetCommits() ([]entities.Commit, error)
	UpdateCommitByEvent(event entities.Event) (*entities.Commit, error)
}

type CommitService struct {
	commitStore persistence.CommitStoreInterface
}

func NewCommitService(commitStore persistence.CommitStoreInterface) *CommitService {
	service := new(CommitService)
	service.commitStore = commitStore
	return service
}

func (s *CommitService) HandleEvent(event entities.Event) (*entities.Commit, error) {
	if event.Type == entities.EVENT_COMMIT {
		// Create commit
		newCommit := entities.Commit{
			PipelineId: event.PipelineId,
			CommitId:   event.CommitId,
			State:      entities.COMMIT_STATE_COMMITTED,
			CommitTime: event.Timestamp,
		}
		return s.CreateCommit(newCommit)

	} else {
		// Update commit
		return s.UpdateCommitByEvent(event)
	}
}

func (s *CommitService) GetCommits() ([]entities.Commit, error) {
	return s.commitStore.GetCommits()
}

func (s *CommitService) CreateCommit(commit entities.Commit) (*entities.Commit, error) {
	// If commit exist for same commit, return error
	_, err := s.commitStore.GetCommitByPipelineIdAndCommitId(commit.PipelineId, commit.CommitId)
	if err == nil {
		// Found
		return nil, errors.New("commit already exist for this commit and pipeline id")
	}
	commit.State = entities.COMMIT_STATE_COMMITTED
	return s.commitStore.CreateCommit(commit)
}

func (s *CommitService) UpdateCommitByEvent(event entities.Event) (*entities.Commit, error) {
	commit, err := s.commitStore.GetCommitByPipelineIdAndCommitId(event.PipelineId, event.CommitId)
	if err != nil {
		return nil, errors.New("no commit matching to update")
	} else {
		if event.Type == entities.EVENT_SUBMIT {
			if event.Status == entities.STATUS_SUCCESS {

				if commit.State == entities.COMMIT_STATE_SUBMITTED {
					// NOOP
					return commit, nil
				}
				if commit.State == entities.COMMIT_STATE_COMMITTED {
					// Submitted
					commit.State = entities.COMMIT_STATE_SUBMITTED
				}
			}
		} else if event.Type == entities.EVENT_APPROVE {
			if event.Status == entities.STATUS_SUCCESS {

				if commit.State == entities.COMMIT_STATE_COMMITTED {
					// Bypass submission
					commit.SubmitTime = commit.CommitTime
				}
				if commit.State == entities.COMMIT_STATE_COMMITTED ||
					commit.State == entities.COMMIT_STATE_SUBMITTED {
					// Approved: update review lead time
					commit.ApprovalTime = event.Timestamp
					commit.State = entities.COMMIT_STATE_APPROVED
					commit.ReviewLeadTime = commit.ApprovalTime.Unix() - commit.SubmitTime.Unix()
					commit.TotalLeadTime = commit.ReviewLeadTime
				}
			}
		} else if event.Type == entities.EVENT_DEPLOY {
			if event.Status == entities.STATUS_SUCCESS {
				if commit.State == entities.COMMIT_STATE_COMMITTED ||
					commit.State == entities.COMMIT_STATE_SUBMITTED {
					// Bypass submission and approval as if they happened at creation time of the commit
					commit.SubmitTime = commit.CommitTime
					commit.ApprovalTime = commit.CommitTime
				}
				// Update state
				commit.State = entities.COMMIT_STATE_DEPLOYED
				// Update DeploymentLeadTime
				commit.DeploymentTime = event.Timestamp
				commit.DeploymentLeadTime = commit.DeploymentTime.Unix() - commit.ApprovalTime.Unix()
				commit.TotalLeadTime = commit.ReviewLeadTime + commit.DeploymentLeadTime
			}
		}
		return s.commitStore.UpdateCommit(*commit)
	}
}
