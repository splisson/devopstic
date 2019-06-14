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
	r := devopstic.BuildEngine(commitHandlers)
	//fmt.Printf("Starting opstic server\n")
	err := r.Run()
	if err != nil {
		panic(err)
	}
}
