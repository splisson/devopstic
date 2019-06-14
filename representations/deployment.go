package representations

type Deployment struct {
	Id          string `json:"id"`
	Timestamp   int64  `json:"timestamp" binding:"required"`   // time of event in nanosecond
	PipelineId  string `json:"pipeline_id" binding:"required"` // key used to group events
	Status      string `json:"status" binding:"required"`      // success, fail, ...
	CommitId    string `json:"commit_id"`                      // Source code commit
	Environment string `json:"environment" binding:"required"` // dev, qa, prod
}

type DeploymentResults struct {
	Items []Deployment `json:"items"`
	Count int          `json:"count"`
	Skip  int          `json:"skip"`
	Limit int          `json:"limit"`
}
