package server

import (
	"net"

	"github.com/turing-era/turingera/log"
	"google.golang.org/grpc"
)

// GrpcConfig grpc服务配置
type GrpcConfig struct {
	Name              string
	Addr              string
	AuthPublicKeyFile string
	RegisterFunc      func(*grpc.Server)
}

// RunGrpcServer 启动grpc服务
func RunGrpcServer(c *GrpcConfig) {
	lis, err := net.Listen("tcp", c.Addr)
	if err != nil {
		panic("grpc cannot listen: " + err.Error())
	}

	var opts []grpc.ServerOption
	// 服务日志拦截器
	opts = append(opts, grpc.UnaryInterceptor(log.ServerLogInterceptor))

	s := grpc.NewServer(opts...)
	c.RegisterFunc(s)

	log.Infof("[%v]server started: %v", c.Name, c.Addr)
	panic(s.Serve(lis))
}
