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
	deploymentStore := persistence.NewDeploymentDBStore(db)
	commitService := services.NewCommitService(commitStore, deploymentStore)
	commitHandlers := handlers.NewCommitHandlers(commitService)
	eventHandlers := handlers.NewEventHandlers(commitService)
	r := devopstic.BuildEngine(commitHandlers, eventHandlers)
	//fmt.Printf("Starting opstic server\n")
	err := r.Run()
	if err != nil {
		panic(err)
	}
}
