package main

import (
	"log"

	"github.com/JaneKetko/Buses/src/client"

	"google.golang.org/grpc"

	"github.com/JaneKetko/Buses/api/proto"
)

func main() {
	conn, err := grpc.Dial(":8001", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	c := proto.NewBusesManagerClient(conn)

	cln := client.NewClient("admin", "admin", c)
	srv := client.NewServer(cln)
	srv.RunServer(":8080")
}
