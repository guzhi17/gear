package main

import (
	"apps/proto/build/go/pb"
	"context"
	"google.golang.org/grpc"
	"log"
	"os"
	"testing"
	"time"
)


const (
	address     = "localhost:50051"
	defaultName = "world"
)
func TestServer_SayHello(t *testing.T) {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.SysEchoQuery{Word: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetWord())
}