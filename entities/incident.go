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
	SourceId       string // Grouping id for incidents from the same source
}

const (
	INCIDENT_STATE_OPEN     = " open"
	INCIDENT_STATE_RESOLVED = "resolved"
)
