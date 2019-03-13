package main

import (
	"log"

	config "github.com/JaneKetko/Buses/src/config"
	routemanager "github.com/JaneKetko/Buses/src/controller"
	dbmanager "github.com/JaneKetko/Buses/src/db"
	"github.com/JaneKetko/Buses/src/server"

	_ "github.com/go-sql-driver/mysql"
)

// func buildContainer() *dig.Container {
// 	container := dig.New()

// 	err := container.Provide(config.GetData)
// 	if err != nil {
// 		log.Fatal("Error with configuration", err)
// 	}

// 	err = container.Provide(dbmanager.Open)
// 	if err != nil {
// 		log.Fatal("Error with opening database", err)
// 	}

// 	err = container.Provide(routemanager.NewRouteManager)
// 	if err != nil {
// 		log.Fatal("Error with creating route manager", err)
// 	}

// 	err = container.Provide(server.NewBusStation)
// 	if err != nil {
// 		log.Fatal("Error with creating busstation", err)
// 	}

// 	return container
// }

func main() {

	config := config.GetData()
	db, err := dbmanager.Open(config)
	if err != nil {
		log.Fatal(err)
	}

	dbmanager := dbmanager.NewDBManager(db)
	routemanager := routemanager.NewRouteManager(dbmanager)
	busstation := server.NewBusStation(routemanager, config)
	busstation.StartServer()
	// container := buildContainer()

	// err := container.Invoke(func(b *server.BusStation) {
	// 	b.StartServer()
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }
}
