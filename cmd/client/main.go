package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"

	pb "github.com/JaneKetko/Buses/api/proto"
)

func main() {
	conn, err := grpc.Dial(":8000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	c := pb.NewBusesManagerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id := 13

	buyreq := pb.IDRequest{
		ID: int64(id),
	}

	buyres, err := c.BuyTicket(ctx, &buyreq)
	if err != nil {
		log.Fatalf("You haven't bought ticket: %v", err)
	}
	log.Printf("You have bought ticket successfully! Your ticket: %v", buyres)
}
