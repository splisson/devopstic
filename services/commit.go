package services

import (
	"errors"
	"github.com/splisson/devopstic/entities"
	"github.com/splisson/devopstic/persistence"
)

type CommitServiceInterface interface {
	CreateCommit(event entities.Commit) (*entities.Commit, error)
	GetCommits() ([]entities.Commit, error)
	UpdateCommitByEvent(commitEvent entities.CommitEvent) (*entities.Commit, error)
}

type CommitService struct {
	commitStore     persistence.CommitStoreInterface
	deploymentStore persistence.DeploymentStoreInterface
}

func NewCommitService(commitStore persistence.CommitStoreInterface, deploymentStore persistence.DeploymentStoreInterface) *CommitService {
	service := new(CommitService)
	service.commitStore = commitStore
	service.deploymentStore = deploymentStore
	return service
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

func (s *CommitService) UpdateCommitByEvent(commitEvent entities.CommitEvent) (*entities.Commit, error) {
	commit, err := s.commitStore.GetCommitByPipelineIdAndCommitId(commitEvent.PipelineId, commitEvent.CommitId)
	if err != nil {
		return nil, errors.New("no commit matching to update")
	} else {
		if commitEvent.Type == entities.COMMIT_EVENT_SUBMIT {
			if commitEvent.Status == entities.STATUS_SUCCESS {

				if commit.State == entities.COMMIT_STATE_SUBMITTED {
					// NOOP
					return commit, nil
				}
				if commit.State == entities.COMMIT_STATE_COMMITTED {
					// Submitted
					commit.State = entities.COMMIT_STATE_SUBMITTED
				}
			}
		} else if commitEvent.Type == entities.COMMIT_EVENT_APPROVE {
			if commitEvent.Status == entities.STATUS_SUCCESS {

				if commit.State == entities.COMMIT_STATE_COMMITTED {
					// Bypass submission
					commit.SubmitTime = commit.CommitTime
				}
				if commit.State == entities.COMMIT_STATE_COMMITTED ||
					commit.State == entities.COMMIT_STATE_SUBMITTED {
					// Approved: update review lead time
					commit.ApprovalTime = commitEvent.Timestamp
					commit.State = entities.COMMIT_STATE_APPROVED
					commit.ReviewLeadTime = commit.ApprovalTime.Unix() - commit.SubmitTime.Unix()
					commit.TotalLeadTime = commit.ReviewLeadTime
				}
			}
		} else if commitEvent.Type == entities.COMMIT_EVENT_DEPLOY {
			// Create deployment
			deployment := entities.Deployment{
				Timestamp:   commitEvent.Timestamp,
				Status:      commitEvent.Status,
				CommitId:    commit.CommitId,
				PipelineId:  commit.PipelineId,
				Environment: commitEvent.Environment,
			}
			if commitEvent.Status == entities.STATUS_SUCCESS {
				if commit.State == entities.COMMIT_STATE_COMMITTED ||
					commit.State == entities.COMMIT_STATE_SUBMITTED {
					// Bypass submission and approval as if they happened at creation time of the commit
					commit.SubmitTime = commit.CommitTime
					commit.ApprovalTime = commit.CommitTime
				}
				// Update state
				commit.State = entities.COMMIT_STATE_DEPLOYED
				// Update DeploymentLeadTime
				commit.DeploymentTime = commitEvent.Timestamp
				commit.DeploymentLeadTime = commit.DeploymentTime.Unix() - commit.ApprovalTime.Unix()
				commit.TotalLeadTime = commit.ReviewLeadTime + commit.DeploymentLeadTime
			}
			_, err := s.deploymentStore.CreateDeployment(deployment)
			if err != nil {
				return nil, errors.New("error creating deployment")
			}

		}
		return s.commitStore.UpdateCommit(*commit)
	}
}
