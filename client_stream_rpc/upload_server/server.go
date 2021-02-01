package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"grpcExample/client_stream_rpc/proto"
	"io"
	"log"
	"net"
)

type UploadServer struct{}

func main() {
	lis, err := net.Listen("tcp", ":6012")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	//构建一个新的服务端对象
	s := grpc.NewServer()
	//向这个服务端对象注册服务
	proto.RegisterUploadServer(s,&UploadServer{})
	//注册服务端反射服务
	reflection.Register(s)

	//启动服务
	s.Serve(lis)

	//可配合ctx实现服务端的动态终止
	//s.Stop()
}

func (*UploadServer) Upload(uploadServer proto.Upload_UploadServer) error {
	for {
		//循环接受客户端传的流数据
		recv, err := uploadServer.Recv()
		//检测到EOF（客户端调用close）
		if err == io.EOF {
			//发送res
			err := uploadServer.SendAndClose(&proto.UploadRes{Msg: "finish"})
			if err != nil {
				return err
			}
			return nil
		} else if err != nil{
			return err
		}
		log.Printf("get a upload data package~ offset:%v, size:%v\n",recv.Offset,recv.Size)
	}
}