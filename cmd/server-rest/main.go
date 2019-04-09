package main

import (
	"log"

	"github.com/JaneKetko/Buses/src/config"
	"github.com/JaneKetko/Buses/src/dbmanager"
	"github.com/JaneKetko/Buses/src/routemanager"
	"github.com/JaneKetko/Buses/src/server"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	cfg := config.GetData()
	db, err := dbmanager.Open(cfg)
	if err != nil {
		log.Fatal(err)
	}

	dbman := dbmanager.NewDBManager(db)
	routeman := routemanager.NewRouteManager(dbman)
	restserv := server.NewBusStation(routeman, cfg)
	restserv.RunServer()
}
