package main

import (
	"log"

	"github.com/JaneKetko/Buses/src/dbmanager"
	"github.com/JaneKetko/Buses/src/grpcserver"
	"github.com/JaneKetko/Buses/src/routemanager"
	"github.com/JaneKetko/Buses/src/server"
	"github.com/JaneKetko/Buses/src/service"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	var sett Config
	err := sett.Parse()
	if err != nil {
		log.Fatalf("Cannot parse settings: %v", err)
	}

	db, err := dbmanager.Open(&dbmanager.DBConfig{
		Login:    sett.Login,
		Passwd:   sett.Passwd,
		Hostname: sett.Hostname,
		Port:     sett.Port,
		DBName:   sett.DBName,
	})
	if err != nil {
		log.Fatal(err)
	}
	dbman := dbmanager.NewDBManager(db)
	routeman := routemanager.NewRouteManager(dbman)
	r := server.NewRESTServer(routeman, sett.PortRESTServer)
	g := grpcserver.NewGRPSServer(routeman, sett.PortGRPCServer)
	srvc := service.NewService(g, r)
	srvc.RunService()
}