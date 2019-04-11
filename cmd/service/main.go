package main

import (
	"log"

	"github.com/JaneKetko/Buses/src/config"
	"github.com/JaneKetko/Buses/src/dbmanager"
	"github.com/JaneKetko/Buses/src/grpcserver"
	"github.com/JaneKetko/Buses/src/routemanager"
	"github.com/JaneKetko/Buses/src/server"
	"github.com/JaneKetko/Buses/src/service"

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
	r := server.NewRESTServer(routeman, cfg)
	g := grpcserver.NewGRPSServer(routeman, cfg)
	srvc := service.NewService(g, r)
	srvc.RunService()
}
