package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"grpcExample/simple_rpc/proto"
	"log"
	"net"
)

type LoginServer struct {}

func main() {
	//监听tcp，6012端口，注意":"
	lis, err := net.Listen("tcp", ":6012")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	//构建一个新的服务端对象
	s := grpc.NewServer()
	//向这个服务端对象注册服务
	proto.RegisterLoginServer(s,&LoginServer{})
	//注册服务端反射服务
	reflection.Register(s)

	//启动服务
	s.Serve(lis)

	//可配合ctx实现服务端的动态终止
	//s.Stop()
}

//判断用户名，密码是否为root,123456，验证正确即返回
func (*LoginServer)Login(ctx context.Context, req *proto.LoginReq) (*proto.LoginRes, error) {
	//为降低复杂度，此处不对ctx进行处理
	if req.Username == "root" && req.Password == "123456" {
		return &proto.LoginRes{Msg: "true"},nil
	} else {
		return &proto.LoginRes{Msg: "false"},nil
	}
}


