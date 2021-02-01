package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"grpcExample/server_stream_rpc/proto"
	"log"
	"net"
)

type DownloadServer struct{}

func main() {
	lis, err := net.Listen("tcp", ":6012")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	//构建一个新的服务端对象
	s := grpc.NewServer()
	//向这个服务端对象注册服务
	proto.RegisterDownloadServer(s,&DownloadServer{})
	//注册服务端反射服务
	reflection.Register(s)

	//启动服务
	s.Serve(lis)

	//可配合ctx实现服务端的动态终止
	//s.Stop()
}

func (*DownloadServer) Download(req *proto.DownloadReq, downloadServer proto.Download_DownloadServer) error {
	offset := req.Offset
	//循环发送数据
	for {
		err := downloadServer.Send(&proto.DownloadRes{
			Offset: offset,
			Size:   4 * 1024,
			Data:   nil,
		})
		if err != nil {
			return err
		}
		offset += 4 * 1024
		if offset >= req.Offset + req.Size {
			break
		}
	}
	return nil
}