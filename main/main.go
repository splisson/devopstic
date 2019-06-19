package main

import (
	"github.com/splisson/devopstic"
	"github.com/splisson/devopstic/handlers"
	"github.com/splisson/devopstic/persistence"
	"github.com/splisson/devopstic/services"
)

func main() {
	db := persistence.NewPostgresqlConnectionWithEnv()
	commitStore := persistence.NewCommitStoreDB(db)
	eventStore := persistence.NewEventStoreDB(db)
	eventService := services.NewEventService(eventStore)
	commitService := services.NewCommitService(commitStore)
	commitHandlers := handlers.NewCommitHandlers(commitService)
	eventHandlers := handlers.NewEventHandlers(eventService, commitService)
	r := devopstic.BuildEngine(commitHandlers, eventHandlers)
	//fmt.Printf("Starting opstic server\n")
	err := r.Run()
	if err != nil {
		panic(err)
	}
}
