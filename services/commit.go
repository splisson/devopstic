package services

import (
	"errors"
	"fmt"
	"github.com/prometheus/common/log"
	"github.com/splisson/devopstic/entities"
	"github.com/splisson/devopstic/persistence"
)

type CommitServiceInterface interface {
	HandleEvent(event entities.Event) (*entities.Commit, *entities.Deployment, error)
	CreateCommit(event entities.Commit) (*entities.Commit, error)
	GetCommits() ([]entities.Commit, error)
	GetCommitByPipelineIdAndCommitId(pipelineId string, commitId string) (*entities.Commit, error)
	UpdateCommitByEvent(commit entities.Commit, event entities.Event) (*entities.Commit, *entities.Deployment, error)
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

func updateLeadTimes(commit *entities.Commit) {
	if commit.ApprovalTime.Unix() >= commit.SubmitTime.Unix() {
		commit.ReviewLeadTime = commit.ApprovalTime.Unix() - commit.SubmitTime.Unix()
	} else {
		commit.ReviewLeadTime = 0
	}
}

func updateDeploymentLeadTimes(commit *entities.Commit, deployment *entities.Deployment) {
	if deployment.Timestamp.Unix() >= commit.ApprovalTime.Unix() {
		deployment.LeadTime = deployment.Timestamp.Unix() - commit.ApprovalTime.Unix()
	} else {
		deployment.LeadTime = 0
	}
}

func (s *CommitService) HandleEvent(event entities.Event) (*entities.Commit, *entities.Deployment, error) {
	var err error = nil
	var commit *entities.Commit
	// PullRequest related
	if event.PullRequestId > 0 {
		commit, err = s.commitStore.GetCommitByPipelineIdAndPullRequestId(event.PipelineId, event.PullRequestId)
	}
	// Find by commitId
	if err != nil || event.PullRequestId <= 0 {
		commit, err = s.commitStore.GetCommitByPipelineIdAndCommitId(event.PipelineId, event.CommitId)
	}
	// Commit not found
	if err != nil {
		// Create new commit
		newCommit := entities.Commit{
			PipelineId:    event.PipelineId,
			CommitId:      event.CommitId,
			CommitTime:    event.Timestamp,
			PullRequestId: event.PullRequestId,
		}
		newCommit.State = entities.COMMIT_STATE_COMMITTED
		commit, err = s.CreateCommit(newCommit)
	}
	// Update new commit
	return s.UpdateCommitByEvent(*commit, event)

	//if err == nil {
	//	// Found => Update
	//	// Update commit
	//	return s.UpdateCommitByEvent(*commit, event)
	//} else {
	//	// Create commit with committed state by default
	//	newCommit := entities.Commit{
	//		PipelineId:    event.PipelineId,
	//		CommitId:      event.CommitId,
	//		CommitTime:    event.Timestamp,
	//		PullRequestId: event.PullRequestId,
	//	}
	//	newCommit.State = entities.COMMIT_STATE_COMMITTED
	//	// Only successful event set time and state on commit
	//	if event.Status == entities.STATUS_SUCCESS {
	//		switch event.Type {
	//		case entities.EVENT_SUBMIT:
	//			newCommit.SubmitTime = event.Timestamp
	//			newCommit.State = entities.COMMIT_STATE_SUBMITTED
	//			break
	//		case entities.EVENT_APPROVE:
	//			newCommit.SubmitTime = event.Timestamp
	//			newCommit.ApprovalTime = event.Timestamp
	//			newCommit.State = entities.COMMIT_STATE_APPROVED
	//			break
	//		case entities.EVENT_DEPLOY:
	//			newCommit.SubmitTime = event.Timestamp
	//			newCommit.ApprovalTime = event.Timestamp
	//			newCommit.State = entities.COMMIT_STATE_DEPLOYED
	//			break
	//		}
	//	}
	//	rCommit, err := s.CreateCommit(newCommit)
	//	return rCommit, nil, err
	//}
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

func (s *CommitService) UpdateCommitByEvent(commit entities.Commit, event entities.Event) (*entities.Commit, *entities.Deployment, error) {
	var rDeployment *entities.Deployment
	var err error

	if event.Type == entities.EVENT_SUBMIT {
		if event.Status == entities.STATUS_SUCCESS {

			if commit.State == entities.COMMIT_STATE_SUBMITTED {
				// NOOP
				return &commit, nil, nil
			}
			if commit.State == entities.COMMIT_STATE_COMMITTED {
				// Submitted
				commit.PullRequestId = event.PullRequestId
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
				// update commit Id with potentially the merge sha from the pull request
				commit.CommitId = event.CommitId
				updateLeadTimes(&commit)
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

			// Deployment
			rDeployment, err = s.createOrUpdateDeployment(commit, event)
			if err != nil {
				return nil, nil, err
			}
		}
	}
	rCommit, err := s.commitStore.UpdateCommit(commit)
	return rCommit, rDeployment, err

}

func (s *CommitService) createOrUpdateDeployment(commit entities.Commit, event entities.Event) (*entities.Deployment, error) {
	var rDeployment *entities.Deployment
	// Is there already a deployment tracked for this commit and environment
	deployment, err := s.deploymentStore.GetDeploymentByCommitIdAndEnvironment(commit.CommitId, event.Environment)
	if err == nil && deployment != nil {
		// Update time
		deployment.Timestamp = event.Timestamp
		updateDeploymentLeadTimes(&commit, deployment)
		rDeployment, err = s.deploymentStore.UpdateDeployment(*deployment)
		if err != nil {
			log.Error(err)
			return nil, errors.New(fmt.Sprintf("error updating deployment for commit %s", commit.CommitId))
		}
	} else {
		// Create deployment
		newDeployment := entities.Deployment{}
		newDeployment.Environment = event.Environment
		newDeployment.CommitId = commit.CommitId
		newDeployment.Status = event.Status
		newDeployment.Timestamp = event.Timestamp
		newDeployment.PipelineId = event.PipelineId
		updateDeploymentLeadTimes(&commit, &newDeployment)
		rDeployment, err = s.deploymentStore.CreateDeployment(newDeployment)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("error creating deployment for commit %s", commit.CommitId))
		}
	}

	return rDeployment, nil
}
