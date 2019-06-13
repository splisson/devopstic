package services

import (
	"github.com/splisson/opstic/persistence"
	"testing"
)

var (
	testEventStore persistence.EventStoreInterface
)
func TestMain(m *testing.M) {
	testEventStore = persistence.NewEventStoreFake()
	m.Run()
}
