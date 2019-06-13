package entities

import (
	"time"
)

type Event struct {
	Model
	Category      string    // build, deploy, etc
	Timestamp     time.Time // time of event in nanosecond
	PipelineId    string    // key used to group events
	Status        string    // success, fail, ...
	Commit        string    // Source code commit
	Environment   string    // dev, qa, prod
	LeadTime      int64     // seconds between build time and deployment
	TimeToRestore int64     // seconds between failure time and recovery
}

const (
	EVENT_CATEGORY_DEPLOY   = "deploy"
	EVENT_CATEGORY_BUILD    = "build"
	EVENT_CATEGORY_INCIDENT = "incident"
	EVENT_STATUS_SUCCESS    = "success"
	EVENT_STATUS_FAILURE    = "failure"
)
