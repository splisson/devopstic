package representations

type Commit struct {
	Id                 string `json:"id"`
	CommitTime         int64  `json:"commit_time"`
	SubmitTime         int64  `json:"submit_time"`
	ApprovalTime       int64  `json:"approval_time"`
	DeploymentTime     int64  `json:"deployment_time"`
	ReviewLeadTime     int64  `json:"review_lead_time"`
	DeploymentLeadTime int64  `json:"deployment_lead_time"`
	TotalLeadTime      int64  `json:"total_lead_time"`
	PipelineId         string `json:"pipeline_id" binding:"required"` // key used to group events
	CommitId           string `json:"commit_id"`                      // Source code commit
}

type CommitResults struct {
	Items []Commit `json:"items"`
	Count int      `json:"count"`
	Skip  int      `json:"skip"`
	Limit int      `json:"limit"`
}
