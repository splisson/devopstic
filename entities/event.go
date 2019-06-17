package entities

import "time"

type Event struct {
	Model
	PipelineId  string
	CommitId    string
	IncidentId  string
	Type        string // Event type
	Status      string // success or failure
	Timestamp   time.Time
	Environment string
}

const (
	EVENT_COMMIT                 = "commit"  // created commit
	EVENT_SUBMIT                 = "submit"  // submit to code review
	EVENT_APPROVE                = "approve" // approve code review
	EVENT_DEPLOY                 = "deploy"
	EVENT_INCIDENT_STATUS_CHANGE = "incident_status_change"
)
