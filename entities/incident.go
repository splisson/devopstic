package entities

import "time"

type Incident struct {
	Model
	OpeningTime    time.Time // time of event
	ResolutionTime time.Time
	PipelineId     string // key used to group events
	State          string // opened, resolved
	Environment    string // dev, qa, prod
	TimeToRestore  int64
	IncidentId     string
}

const (
	INCIDENT_STATE_OPENED   = " opened"
	INCIDENT_STATE_RESOLVED = "resolved"
)
