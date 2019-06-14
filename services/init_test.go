package services

import (
	"github.com/splisson/devopstic/persistence"
	"testing"
)

var (
	testDeploymentStore persistence.DeploymentStoreInterface
	testCommitStore     persistence.CommitStoreInterface
	testIncidentStore   persistence.IncidentStoreInterface
)

func TestMain(m *testing.M) {
	testDeploymentStore = persistence.NewDeploymentStoreFake()
	testCommitStore = persistence.NewCommitStoreFake()
	testIncidentStore = persistence.NewIncidentStoreFake()
	m.Run()
}
