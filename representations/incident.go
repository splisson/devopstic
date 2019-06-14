package representations

type Incident struct {
	Id            string `json:"id"`
	Timestamp     int64  `json:"timestamp" binding:"required"`   // time of event in nanosecond
	PipelineId    string `json:"pipeline_id" binding:"required"` // key used to group events
	Status        string `json:"status" binding:"required"`      // success, fail, ...
	Environment   string `json:"environment" binding:"required"` // dev, qa, prod
	TimeToRestore int64  `json:"time_to_restore"`
	IncidentId    string `json:"incident_id"`
}

type IncidentResults struct {
	Items []Incident `json:"items"`
	Count int        `json:"count"`
	Skip  int        `json:"skip"`
	Limit int        `json:"limit"`
}
