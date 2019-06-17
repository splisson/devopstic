package services

import (
	"github.com/splisson/devopstic/persistence"
	"testing"
)

var (
	testEventStore    persistence.EventStoreInterface
	testCommitStore   persistence.CommitStoreInterface
	testIncidentStore persistence.IncidentStoreInterface
)

func TestMain(m *testing.M) {
	testEventStore = persistence.NewEventStoreFake()
	testCommitStore = persistence.NewCommitStoreFake()
	testIncidentStore = persistence.NewIncidentStoreFake()
	m.Run()
}
