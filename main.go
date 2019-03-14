package main

import (
	"log"

	config "github.com/JaneKetko/Buses/src/config"
	routemanager "github.com/JaneKetko/Buses/src/controller"
	dbmanager "github.com/JaneKetko/Buses/src/db"
	"github.com/JaneKetko/Buses/src/server"

	_ "github.com/go-sql-driver/mysql"
)

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
}
