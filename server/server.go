package main


import (
	"file-transport/proto"
	"file-transport/transport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)



func main() {
	listener, err := net.Listen("tcp", ":4040")
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


