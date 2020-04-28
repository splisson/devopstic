package representations

type PagerDutyIncident struct {
	Id       string    `json:"id"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Event    string     `json:"event"`
	Incident PDIncident `json:"incident"`
}

const (
	PAGERDUTY_INCIDENT_STATUS_TRIGGERED = "triggered"
)

type PDIncident struct {
	IncidentNumber int64  `json:"incident_number"`
	Status         string `json:"status"`
}
