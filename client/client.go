package main

import (
	"context"
	"file-transport/proto"
	"google.golang.org/grpc"
	"log"
	"os"
	"path"
)

func main() {
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
