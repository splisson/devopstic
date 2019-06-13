package representations

import "time"

type Event struct {
	Category      string    `json:"category" binding:"required"`    // build, deploy, etc
	Timestamp     time.Time `json:"timestamp" binding:"required"`   // time of event in nanosecond
	PipelineId    string    `json:"pipeline_id" binding:"required"` // key used to group events
	Status        string    `json:"status" binding:"required"`      // success, fail, ...
	Commit        string    `json:"commit" binding:"required"`      // Source code commit
	Environment   string    `json:"environment" binding:"required"` // dev, qa, prod
	LeadTime      string    `json:"leadtime"`
	TimeToRestore string    `json:"time_to_restore"`
}
