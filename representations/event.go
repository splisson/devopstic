package representations

type Event struct {
	Id          string `json:"id"`
	PipelineId  string `json:"pipeline_id"`
	CommitId    string `json:"commit_id"`
	IncidentId  string `json:"incident_id"`
	Timestamp   int64  `json:"timestamp"`
	Type        string `json:"type"`   // Event type
	Status      string `json:"status"` // success or failure
	Environment string `json:"environment"`
}

type EventResults struct {
	Items []Event `json:"items"`
	Count int     `json:"count"`
	Skip  int     `json:"skip"`
	Limit int     `json:"limit"`
}
