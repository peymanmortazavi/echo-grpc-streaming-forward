package main

import (
	"fmt"
	"time"
	"webstreaming/webstreaming/gen"
)

type rpcStreaming struct {
	gen.UnimplementedServiceServer
}

func (r *rpcStreaming) Process(request *gen.ProcessRequest, server gen.Service_ProcessServer) error {
	for i := 0; i < int(request.Count); i++ {
		server.Send(&gen.ProcessResponse{
			Text: fmt.Sprintf("Iteration %d", i),
		})
		time.Sleep(500 * time.Millisecond)
	}

	return nil
}
