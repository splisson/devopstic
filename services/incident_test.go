package services

import (
	"github.com/google/uuid"
	"github.com/splisson/devopstic/entities"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var (
	testIncident = entities.Incident{
		OpeningTime: time.Now(),
		IncidentId:  uuid.New().String(),
		PipelineId:  "test",
		Status:      "failure",
		Environment: "unit_test",
	}
)

func TestIncidentRecovery(t *testing.T) {

	incidentService := NewIncidentService(testIncidentStore)

	t.Run("should fill time to restore when existing incident in failure status", func(t *testing.T) {
		newIncident := testIncident
		newIncident.OpeningTime = time.Now().Add(-5 * time.Minute)
		incident, err := incidentService.CreateOrUpdateIncident(newIncident)
		assert.Nil(t, err, "no error")
		assert.NotNil(t, incident, "incident not nil")
		newIncident = testIncident
		newIncident.Status = entities.STATUS_SUCCESS
		newIncident.OpeningTime = time.Now()
		currentId := incident.ID
		incident, err = incidentService.CreateOrUpdateIncident(newIncident)
		assert.Nil(t, err, "no error")
		assert.NotNil(t, incident, "incident not nil")
		assert.Equal(t, currentId, incident.ID, "should have same ID (update)")
		assert.True(t, incident.TimeToRestore > 0, "should have time to restore > 0")
	})

}
