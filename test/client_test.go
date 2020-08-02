package main

import (
	"context"
	"file-transport/proto"
	"file-transport/transport"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"path"
	"testing"
)

const (
	port = 4040
	chunkSize = 2 << 20
)

func createServer()  {
	address := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", address)
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

func createClient() proto.FileTransferServiceClient {
	target := fmt.Sprintf("localhost:%d", port)
	conn, _ := grpc.Dial(target, grpc.WithInsecure())
	client := proto.NewFileTransferServiceClient(conn)
	return client
}

func TestLaunchServer(t *testing.T)  {
	createServer()
}

func TestPutFileIfNotExists(t *testing.T) {
	//go createServer()
	conn, err := grpc.Dial("localhost:4040", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	client := proto.NewFileTransferServiceClient(conn)

	filePath := "/home/mark/Downloads"
	fileName := "jdk-8u261-linux-x64.tar.gz"

	fullPath := path.Join(filePath, fileName)
	fileInfo, err := os.Stat(fullPath)
	if err != nil {
		log.Fatalln("file stat error", fullPath, err)
	}

	dstPath := "/home/mark"
	fileMetaData := proto.FileMeta{
		Filename: fileName,
		Path:     dstPath,
		Md5:      "1522f0c6380fa4993e932eb0c6006ef1",
		Sha1:     "649332adb37dc23f3c48f81c5b9210f207e4e3e2",
		FileSize: fileInfo.Size(),
		Offset:   0,
	}

	req := &proto.FilePutRequest{
		FileMeta: &fileMetaData,
	}
	response, err := client.Put(context.TODO(), req)
	log.Print(response, err)
}

func TestPutFile(t *testing.T) {
	//go createServer()
	conn, err := grpc.Dial("localhost:4040", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	client := proto.NewFileTransferServiceClient(conn)

	filePath := "/home/mark/Downloads"
	fileName := "jdk-8u261-linux-x64.tar.gz"

	fullPath := path.Join(filePath, fileName)
	fileInfo, err := os.Stat(fullPath)
	if err != nil {
		log.Fatalln("file stat error", fullPath, err)
	}

	dstPath := "/home/mark"
	fileMetaData := proto.FileMeta{
		Filename: fileName,
		Path:     dstPath,
		Md5:      "1522f0c6380fa4993e932eb0c6006ef1",
		Sha1:     "649332adb37dc23f3c48f81c5b9210f207e4e3e2",
		FileSize: fileInfo.Size(),
		Offset:   0,
	}

	req := &proto.FilePutRequest{
		FileMeta: &fileMetaData,
	}
	response, err := client.Put(context.TODO(), req)
	log.Print(response, err)
}

func TestGetFile (t *testing.T)  {
	//go createServer()
	//C:\Users\Administrator\Downloads\WXWork_3.0.26.2606.exe
	filename := "WXWork_3.0.26.2606.exe"
	filepath := "C:\\Users\\Administrator\\Downloads\\"
	fileMeta := &proto.FileMeta{
		Filename: filename,
		Path: filepath,
		ChunkSize: 1 << 20,
	}

	req := &proto.FileGetRequest{
		FileMeta: fileMeta,
	}

	client := createClient()

	resp, err := client.Get(context.TODO(), req)
	fmt.Println(resp, err)
}



