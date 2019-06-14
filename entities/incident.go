package entities

import "time"

type Incident struct {
	Model
	OpeningTime    time.Time // time of event
	ResolutionTime time.Time
	PipelineId     string // key used to group events
	Status         string // success, fail, ...
	Environment    string // dev, qa, prod
	TimeToRestore  int64
	IncidentId     string
}
