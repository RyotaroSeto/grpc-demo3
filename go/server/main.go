package main

import (
	"fmt"
	"log"
	"net"

	"grpc-demo3/deepthought"

	"google.golang.org/grpc"
)

const portNumber = 13333

func main() {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", portNumber))
	if err != nil {
		log.Fatalf("cannot create listener: %s", err)
	}
	serv := grpc.NewServer()
	deepthought.RegisterComputeServer(serv, &Server{})
	err = serv.Serve(l)
	if err != nil {
		log.Fatalf("impossible to serve: %s", err)
	}
}
