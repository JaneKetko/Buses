package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

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

	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM, syscall.SIGINT)

	db, err := dbmanager.Open(&dbmanager.DBConfig{
		Login:   sett.Login,
		Passwd:  sett.Passwd,
		Address: sett.Address,
		DBName:  sett.DBName,
	})
	if err != nil {
		log.Fatal(err)
	}
	dbman := dbmanager.NewDBManager(db)
	routeman := routemanager.NewRouteManager(dbman)
	r := server.NewRESTServer(routeman, sett.PortRESTServer)
	g := grpcserver.NewGRPCServer(routeman, sett.PortGRPCServer)
	srvc := service.NewService(g, r)

	go func() {
		sig := <-gracefulStop
		fmt.Printf("caught sig: %+v", sig)
		srvc.StopService()
		os.Exit(0)
	}()
	srvc.RunService()
}
