package services

import (
	"github.com/splisson/devopstic/persistence"
	"os"
	"testing"
)

var (
	testEventStore      persistence.EventStoreInterface
	testCommitStore     persistence.CommitStoreInterface
	testIncidentStore   persistence.IncidentStoreInterface
	testDeploymentStore persistence.DeploymentStoreInterface
)

func TestMain(m *testing.M) {
	if os.Getenv("TEST_WITH_POSTGRES") == "true" {
		db := persistence.NewPostgresqlConnectionLocalhost()
		testEventStore = persistence.NewEventStoreDB(db)
		testCommitStore = persistence.NewCommitStoreDB(db)
		testDeploymentStore = persistence.NewDeploymentStoreDB(db)
		testIncidentStore = persistence.NewIncidentStoreDB(db)
		persistence.CreateTables(db)
	} else {
		testEventStore = persistence.NewEventStoreFake()
		testCommitStore = persistence.NewCommitStoreFake()
		testIncidentStore = persistence.NewIncidentStoreFake()
		testDeploymentStore = persistence.NewDeploymentStoreFake()
	}

	m.Run()
}
