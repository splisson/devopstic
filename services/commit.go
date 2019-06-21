package services

import (
	"errors"
	"github.com/splisson/devopstic/entities"
	"github.com/splisson/devopstic/persistence"
)

type CommitServiceInterface interface {
	HandleEvent(event entities.Event) (*entities.Commit, error)
	CreateCommit(event entities.Commit) (*entities.Commit, error)
	GetCommits() ([]entities.Commit, error)
	GetCommitByPipelineIdAndCommitId(pipelineId string, commitId string) (*entities.Commit, error)
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

func updateLeadTimes(commit *entities.Commit) {
	if commit.ApprovalTime.Unix() >= commit.SubmitTime.Unix() {
		commit.ReviewLeadTime = commit.ApprovalTime.Unix() - commit.SubmitTime.Unix()
	} else {
		commit.ReviewLeadTime = 0
	}
	if commit.DeploymentTime.Unix() >= commit.ApprovalTime.Unix() {
		commit.DeploymentLeadTime = commit.DeploymentTime.Unix() - commit.ApprovalTime.Unix()
	} else {
		commit.DeploymentLeadTime = 0
	}
	commit.TotalLeadTime = commit.ReviewLeadTime + commit.DeploymentLeadTime
}

func (s *CommitService) HandleEvent(event entities.Event) (*entities.Commit, error) {
	// If commit exist for same commit, return error
	_, err := s.commitStore.GetCommitByPipelineIdAndCommitId(event.PipelineId, event.CommitId)
	if err == nil {
		// Found => Update
		// Update commit
		return s.UpdateCommitByEvent(event)
	} else {
		// Create commit with committed state by default
		newCommit := entities.Commit{
			PipelineId: event.PipelineId,
			CommitId:   event.CommitId,
			CommitTime: event.Timestamp,
		}
		newCommit.State = entities.COMMIT_STATE_COMMITTED
		// Only successful event set time and state on commit
		if event.Status == entities.STATUS_SUCCESS {
			switch event.Type {
			case entities.EVENT_SUBMIT:
				newCommit.SubmitTime = event.Timestamp
				newCommit.State = entities.COMMIT_STATE_SUBMITTED
				break
			case entities.EVENT_APPROVE:
				newCommit.SubmitTime = event.Timestamp
				newCommit.ApprovalTime = event.Timestamp
				newCommit.State = entities.COMMIT_STATE_APPROVED
				break
			case entities.EVENT_DEPLOY:
				newCommit.SubmitTime = event.Timestamp
				newCommit.ApprovalTime = event.Timestamp
				newCommit.DeploymentTime = event.Timestamp
				newCommit.State = entities.COMMIT_STATE_DEPLOYED
				break
			}
		}
		return s.CreateCommit(newCommit)
	}
}

func (s *CommitService) GetCommits() ([]entities.Commit, error) {
	return s.commitStore.GetCommits()
}

func (s *CommitService) GetCommitByPipelineIdAndCommitId(pipelineId string, commitId string) (*entities.Commit, error) {
	return s.commitStore.GetCommitByPipelineIdAndCommitId(pipelineId, commitId)
}

func (s *CommitService) CreateCommit(commit entities.Commit) (*entities.Commit, error) {
	_, err := s.commitStore.GetCommitByPipelineIdAndCommitId(commit.PipelineId, commit.CommitId)
	if err == nil {
		return nil, errors.New("commit with same pipeline_id and commit_id already exists")
	}
	if commit.State == "" {
		commit.State = entities.COMMIT_STATE_COMMITTED
	}
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
					commit.SubmitTime = event.Timestamp
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
					updateLeadTimes(commit)
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
				updateLeadTimes(commit)
			}
		}
		return s.commitStore.UpdateCommit(*commit)
	}
}
