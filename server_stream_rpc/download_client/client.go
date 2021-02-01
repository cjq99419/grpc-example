package main

import (
	"context"
	"google.golang.org/grpc"
	"grpcExample/server_stream_rpc/proto"
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
	client := proto.NewDownloadClient(grpcConn)

	//设置ctx超时（根据情况设定）
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//和简单rpc不同，此时获得的不是res，而是一个client的对象，通过这个连接对象去读取数据
	downloadClient,err := client.Download(ctx,&proto.DownloadReq{
		Path:   "../test.txt",
		Offset: 0,
		Size:   64 * 1024,
	})
	if err != nil {
		log.Fatalln(err)
	}

	//循环处理数据，当监测到读取完成后退出
	for {
		res, err := downloadClient.Recv()
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("get a date package~ offset:%v, size:%v\n",res.Offset,res.Size)
		if res.Size + res.Offset >= 64 * 1024 {
			break
		}
	}

	log.Println("download over~")
}