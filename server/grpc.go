package server

import (
	"net"

	"google.golang.org/grpc"

	"github.com/turing-era/turingera-shared/auth"
	"github.com/turing-era/turingera-shared/log"
)

const (
	AuthPublicKeyTypeRSA256 = 0
	AuthPublicKeyTypeES256 = 1
)

// GrpcConfig grpc服务配置
type GrpcConfig struct {
	Name              string
	Addr              string
	AuthPublicKeyFile string

	RegisterFuncs     []func(*grpc.Server)
}

// RunGrpcServer 启动grpc服务
func RunGrpcServer(c *GrpcConfig) {
	lis, err := net.Listen("tcp", c.Addr)
	if err != nil {
		panic("grpc cannot listen: " + err.Error())
	}
	var interceptors []grpc.UnaryServerInterceptor
	// 日志拦截器
	interceptors = append(interceptors, log.ServerLogInterceptor)
	// 鉴权拦截器
	log.Debugf("serverConfig: %+v", c)
	if len(c.AuthPublicKeyFile) > 0 {
		in, err := auth.LoadInterceptor(c.AuthPublicKeyFile)
		if err != nil {
			panic("cannot create auth intercept: " + err.Error())
		}
		interceptors = append(interceptors, in)
	}
	s := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptors...))
	for _, f := range c.RegisterFuncs {
		f(s)
	}
	log.Infof("[%v]server started: %v", c.Name, c.Addr)
	panic(s.Serve(lis))
}
