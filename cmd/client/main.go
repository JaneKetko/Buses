package main

import (
	"log"

	"google.golang.org/grpc"

	"github.com/JaneKetko/Buses/api/proto"
	cl "github.com/JaneKetko/Buses/src/client/workserver"
)

func main() {
	conn, err := grpc.Dial(":8001", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	c := proto.NewBusesManagerClient(conn)

	cln := cl.NewClient("admin", "admin", c)
	srv := cl.NewServer(cln)
	srv.RunServer(":8080")
}
