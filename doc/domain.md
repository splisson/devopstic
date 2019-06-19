# Domain
## Event
An event is sent to signal a transition of state for a commit, or an incident.
Deployments are a type of event.
## Commit
A commit is the base unit we track in order to compute the metrics related to code review and deployments.
### Lifecycle
![commit state chart](commit_state_chart.png "Commit state chart")
## Incident
An incident is created when the incident is opened and updated when resolved.
### Lifecycle
![incident state chart](incident_state_chart.png "Incident state chart")
