package transport

import (
	"context"
	"file-transport/proto"
	"file-transport/util"
	"log"
	"os"
	"path"
	"path/filepath"
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
	FilePutChannel      chan CreateFileAction
}

func (peer *Peer) CreateFileIfNotExists() {
	var action CreateFileAction
	for {
		select {
		case action = <-peer.FilePutChannel:
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

	peer.FilePutChannel <- action

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
	absolutePath := filepath.Join(fileMeta.Path, fileMeta.Filename)
	fileInfo, err := os.Stat(absolutePath)
	fileGetResponse := &proto.FileGetResponse{
		FileMeta: fileMeta,
	}
	//文件打开有问题 比如文件不存在 拼写错误啥的
	if err != nil {
		return fileGetResponse, err
	}
	md5, err := util.GetFileMd5(absolutePath)
	sha1,  err:= util.GetFileSha1(absolutePath)

	fileMeta.Md5 = md5
	fileMeta.Sha1 = sha1
	fileMeta.Md5 = md5
	fileMeta.Sha1 = sha1
	fileMeta.Offset = 0
	fileMeta.FileSize = fileInfo.Size()

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
	absolutePath := filepath.Join(fileMeta.Path, fileMeta.Filename)
	_, err := os.Stat(absolutePath)

	if err != nil {
		log.Print("file read error", err)
	}

	file, err := os.Open(absolutePath)

	response := &proto.ChunkGetResponse{
		FileMeta: fileMeta,
	}
	if err != nil {
		return response, err
	}
	//log.Fatal("chunkSize offset {} {}", fileMeta.ChunkSize, fileMeta.Offset)
	content := make([]byte, fileMeta.ChunkSize)
	_, err = file.ReadAt(content, fileMeta.Offset)
	defer file.Close()
	response.Content = content

	return response, nil
}
