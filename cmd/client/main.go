package main

import (
	"fmt"
	"log"

	"google.golang.org/grpc"

	"github.com/JaneKetko/Buses/api/proto"
	cl "github.com/JaneKetko/Buses/src/client/workserver"
)

func main() {
	sett := new(Config)
	err := sett.Parse()
	if err != nil {
		log.Fatalf("Cannot parse settings: %v", err)
	}
	conn, err := grpc.Dial(sett.PortGRPCServer, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	c := proto.NewBusesManagerClient(conn)

	var name string
	fmt.Println("Enter your name:")
	fmt.Scan(&name)

	cln := cl.NewClient(name, "user", c)
	srv := cl.NewServer(cln)
	srv.RunServer(sett.PortClient)
}
