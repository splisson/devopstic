package entities

import (
	"time"
)

type Commit struct {
	Model
	State              string
	CommitTime         time.Time
	SubmitTime         time.Time
	ApprovalTime       time.Time
	DeploymentTime     time.Time
	ReviewLeadTime     int64
	DeploymentLeadTime int64
	TotalLeadTime      int64  // ReviewLeadTime + DeploymentLeadTime
	PipelineId         string // external unique id of the component/pipeline tracked
	CommitId           string // external unique id of the commit
}

const (
	COMMIT_STATE_COMMITTED = "committed"
	COMMIT_STATE_SUBMITTED = "submitted"
	COMMIT_STATE_APPROVED  = "approved"
	COMMIT_STATE_DEPLOYED  = "deployed"
)

const (
	COMMIT_EVENT_COMMIT  = "commit"  // created commit
	COMMIT_EVENT_SUBMIT  = "submit"  // submit to code review
	COMMIT_EVENT_APPROVE = "approve" // approve code review
	COMMIT_EVENT_DEPLOY  = "deploy"
)

const (
	COMMIT_EVENT_SUCCESS = "success"
	COMMIT_EVENT_FAILURE = "failure"
)

// Event related to a commit identified by pipelineID and commitId
type CommitEvent struct {
	PipelineId  string
	CommitId    string
	Type        string // Event type
	Status      string // success or failure
	Timestamp   time.Time
	Environment string
}
