package main

import (
	"log"
	"webstreaming/webstreaming/gen"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	connection, err := grpc.Dial(
		"localhost:40000",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock())

	if err != nil {
		log.Fatalf("failed to connect: %s", err)
	}

	client := gen.NewServiceClient(connection)
	c, err := client.Process(context.Background(), &gen.ProcessRequest{Count: 10})
	if err != nil {
		log.Fatalf("failed to connect")
	}

	for {
		msg, err := c.Recv()
		if err != nil {
			log.Fatalf("failed: %s", err)
		}
		log.Println(msg.Text)
	}

}
