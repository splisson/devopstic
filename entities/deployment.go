package entities

import (
	"time"
)

type Deployment struct {
	Model
	Timestamp   time.Time // time of event
	PipelineId  string    // key used to group events
	Status      string    // success, fail, ...
	CommitId    string    // Source code commit
	Environment string    // dev, qa, prod
}
