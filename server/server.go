package main

import (
	"file-transport/proto"
	"file-transport/transport"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

var (
	h bool
)


func main() {

	port := flag.Int("port", 4040, "grpc server port")
	network := flag.String("network", "tcp", "network type")


	flag.Parse()

	address := fmt.Sprintf(":%d", *port)
	listener, err := net.Listen(*network, address)
	if err != nil {
		panic(err)
	}

	srv := grpc.NewServer()

	peer := &transport.Peer{
		FilePutChannel: make(chan transport.CreateFileAction, 10000),
	}
	go peer.CreateFileIfNotExists()
	proto.RegisterFileTransferServiceServer(srv, peer)
	reflection.Register(srv)

	if e := srv.Serve(listener); e != nil {
		panic(e)
	}

}


