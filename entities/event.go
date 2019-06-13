package entities

import (
	"time"
)

type Event struct {
	Model
	Category		string		// build, deploy, etc
	Timestamp 		time.Time   // time of event in nanosecond
	PipelineId		string		// key used to group events
	Status			string		// success, fail, ...
	Commit			string		// Source code commit
	Environment     string      // dev, qa, prod
	LeadTime		int64       // seconds between build time and deployment
}

const (
	EVENT_CATEGORY_DEPLOY = "deploy"
	EVENT_CATEGORY_BUILD  = "build"
)
