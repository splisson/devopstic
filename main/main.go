package main

import (
	"github.com/splisson/opstic"
	"github.com/splisson/opstic/handlers"
	"github.com/splisson/opstic/persistence"
	"github.com/splisson/opstic/services"
)

func main() {
	db := persistence.NewPostgresqlConnectionWithEnv()
	eventStore := persistence.NewEventDBStore(db)
	eventService := services.NewEventService(eventStore)
	eventHandlers := handlers.NewEventHandlers(eventService)
	r := opstic.BuildEngine(eventHandlers)
	//fmt.Printf("Starting opstic server\n")
	err := r.Run()
	if err != nil {
		panic(err)
	}
}
