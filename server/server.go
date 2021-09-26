package server

import (
	"context"
	"fmt"
	v1 "go_grpc/api/protocol/v1"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
)

func RunServer(ctx context.Context, v1API v1.ToDoServerServer, port string) error {
	listen, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer()
	v1.RegisterToDoServerServer(server, v1API)
	c := make(chan os.Signal, 1)
	go func() {
		for range c {
			log.Println("shut down grpc server")
			server.GracefulStop()
			<-ctx.Done()
		}
	}()
	fmt.Println("grpc server start ...")
	return server.Serve(listen)
}
