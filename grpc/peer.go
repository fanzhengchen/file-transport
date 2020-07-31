package grpc

import (
	"context"
	"file-transport/proto"
	"log"
	"os"
	"path"
	"sync"
)

type CreateFileAction struct {
	fileMeta *proto.FileMeta
	// 出现结果以后接可以通知接口返回了 暂时不太熟悉 go java里的话，可以注册回调函数或者 lockSupport主动唤醒线程 通过线程安全queue通知事件也是可以的
	done chan bool
}

type ReadFileAction struct {
	fileMeta *proto.FileMeta
	done chan bool
}


type Peer struct {
	fileRWMutex 		sync.RWMutex
	filePutChannel      chan CreateFileAction
	fileReadConcurrency int32
	bufferSize			int32
}

func (peer *Peer) createFileIfNotExists() {
	var action CreateFileAction
	for {
		select {
		case action = <-peer.filePutChannel:
			fileMetaData := action.fileMeta
			fileName := fileMetaData.Filename
			filePath := fileMetaData.Path
			absolutePath := path.Join(filePath, fileName)
			_, err := os.Stat(absolutePath)
			if err != nil {
				file, err := os.Create(absolutePath)
				if err != nil {
					action.done <- false
				} else {
					err = file.Truncate(action.fileMeta.FileSize)
					if err != nil {
						action.done <- false
					} else {
						action.done <- true
					}
				}
			}
		}
	}
}

func (peer *Peer) Put(ctx context.Context, request *proto.FilePutRequest) (*proto.FilePutResponse, error) {
	action := CreateFileAction{
		fileMeta: request.FileMeta,
		done: make(chan bool, 1),
	}

	peer.filePutChannel <- action

	//阻塞直到返回 go routine类似于 thread 阻塞的时候我们让出cpu, 切换不会导致 futex api(linux)之类的系统调用
	ack := <- action.done
	response := &proto.FilePutResponse{
		FileMeta: request.FileMeta,
		Ack:      ack,
	}
	return response, nil
}

func (peer *Peer) Get(ctx context.Context, request *proto.FileGetRequest) (*proto.FileGetResponse, error) {
	fileMeta := request.FileMeta
	absolutePath := path.Join(fileMeta.Path, fileMeta.Filename)
	_, err := os.Stat(absolutePath)
	fileGetResponse := &proto.FileGetResponse{
		FileMeta: fileMeta,
	}
	if err != nil {
		return fileGetResponse, err
	}
	return fileGetResponse, nil
}

func (peer *Peer) PutChunk(ctx context.Context, request *proto.ChunkPutRequest) (*proto.ChunkPutResponse, error) {
	fileMeta := request.FileMeta
	absolutePath := path.Join(fileMeta.Path, fileMeta.Filename)
	_, err := os.Stat(absolutePath)

	if err != nil {
		log.Print("file read error", err)
	}

	file, err := os.Open(absolutePath)
	file.Seek(fileMeta.Offset, 0)

	content := request.Content
	_, err = file.Write(content)

	response := &proto.ChunkPutResponse{
		FileMeta: fileMeta,
		Ack: true,
	}
	if err != nil {

	}
	return response, nil
}

func (peer *Peer) GetChunk(ctx context.Context, request *proto.ChunkGetRequest) (*proto.ChunkGetResponse, error) {
	fileMeta := request.FileMeta
	absolutePath := path.Join(fileMeta.Path, fileMeta.Filename)
	_, err := os.Stat(absolutePath)

	if err != nil {
		log.Print("file read error", err)
	}

	file, err := os.Open(absolutePath)
	file.Seek(fileMeta.Offset, 0)

	content := make([]byte, fileMeta.ChunkSize)
	_, err = file.Write(content)

	response := &proto.ChunkGetResponse{
		FileMeta: fileMeta,
		Content: content,
	}
	if err != nil {
		return response, err
	}
	return response, nil
}
