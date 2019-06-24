package main

import (
	"github.com/splisson/devopstic"
	"github.com/splisson/devopstic/handlers"
	"github.com/splisson/devopstic/persistence"
	"github.com/splisson/devopstic/services"
)

func main() {
	db := persistence.NewPostgresqlConnectionWithEnv()
	persistence.CreateTables(db)
	commitStore := persistence.NewCommitStoreDB(db)
	eventStore := persistence.NewEventStoreDB(db)
	incidentStore := persistence.NewIncidentStoreDB(db)
	eventService := services.NewEventService(eventStore)
	commitService := services.NewCommitService(commitStore)
	incidentService := services.NewIncidentService(incidentStore)
	commitHandlers := handlers.NewCommitHandlers(commitService)
	eventHandlers := handlers.NewEventHandlers(eventService, commitService, incidentService)
	githubEventHandlers := handlers.NewGithubEventHandlers(eventService, commitService, incidentService)
	r := devopstic.BuildEngine(commitHandlers, eventHandlers, githubEventHandlers)
	//fmt.Printf("Starting opstic server\n")
	err := r.Run()
	if err != nil {
		panic(err)
	}
}
