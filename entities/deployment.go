package entities

import "time"

type Deployment struct {
	Model
	Status      string
	Timestamp   time.Time
	LeadTime    int64
	CommitId    string
	PipelineId  string
	Environment string
}
