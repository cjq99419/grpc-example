package main

import (
	"context"
	"google.golang.org/grpc"
	"grpcExample/client_stream_rpc/proto"
	"log"
	"time"
)

func main(){
	//创立grpc连接
	grpcConn, err := grpc.Dial("127.0.0.1"+":6012", grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}

	//通过grpc连接创建一个客户端实例对象
	client := proto.NewUploadClient(grpcConn)

	//设置ctx超时（根据情况设定）
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//和简单rpc不同，此时获得的不是res，而是一个client的对象，通过这个连接对象去读取数据
	uploadClient,err := client.Upload(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	var offset int64
	var size int64
	size = 4 * 1024

	//循环处理数据，当大于64kb退出
	for {
		err := uploadClient.Send(&proto.UploadReq{
			Path:   "../test.txt",
			Offset: offset,
			Size:   size,
			Data:   nil,
		})
		if err != nil {
			log.Fatalln(err)
		}
		offset += size
		//发送超过64KB，调用CloseAndRecv方法接收response
		if offset >= 64 * 1024 {
			res, err := uploadClient.CloseAndRecv()
			if err != nil {
				log.Fatalln(err)
			}
			log.Println("upload over~, response is ",res.Msg)
			break
		}
	}
}