package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"webstreaming/webstreaming/gen"

	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

func proxyWebSocket(client gen.ServiceClient) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {
		var globalErr error

		handler := websocket.Handler(func(ws *websocket.Conn) {
			defer ws.Close()

			msg := ""
			if err := websocket.Message.Receive(ws, &msg); err != nil {
				globalErr = fmt.Errorf("error receiving ws: %s", err)
				return
			}

			initialRequest := &gen.ProcessRequest{}
			if err := protojson.Unmarshal([]byte(msg), initialRequest); err != nil {
				globalErr = fmt.Errorf("error unmarshaling ws: %s", err)
				return
			}

			proxyClient, err := client.Process(c.Request().Context(), initialRequest)
			if err != nil {
				globalErr = fmt.Errorf("failed to proxy to client: %s", err)
				return
			}

			// start the response loop.
			for {
				response, err := proxyClient.Recv()

				if err == io.EOF {
					log.Printf("received EOF from gRPC server")
					return
				}

				if err != nil {
					log.Printf("failed to read from grpc server: %s", err)
					continue
				}

				jsonResult, err := protojson.Marshal(response)
				if err != nil {
					log.Printf("failed to json marshal the response form grpc server: %s", err)
					continue
				}
				if _, err := ws.Write(jsonResult); err != nil {
					log.Printf("failed to forward response from grpc server: %s", err)
				}
			}

		})

		handler.ServeHTTP(c.Response(), c.Request())
		return globalErr
	})
}

func main() {
	selfConnection, err := grpc.Dial(
		"localhost:40000",
		grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("failed to connect: %s", err)
	}

	e := echo.New()
	e.GET("/ws", proxyWebSocket(gen.NewServiceClient(selfConnection)))
	e.Static("/", "./templates/web")

	httpServer := &http.Server{Addr: "localhost:5000", Handler: e}
	go func() {
		log.Printf("starting http server at :5000")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http server failed: %s", err)
		}
	}()

	server := grpc.NewServer()
	gen.RegisterServiceServer(server, &rpcStreaming{})
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 40000))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("starting grpc server at :40000")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("grpc server failed: %s", err)
	}
}
