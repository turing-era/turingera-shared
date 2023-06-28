package server

import (
	"context"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/thinkeridea/go-extend/exnet"
	"github.com/turing-era/turingera-shared/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

// GrpcGatewayConfig 网关服务子配置
type GrpcGatewayConfig struct {
	ServerAddr    string
	RegisterFuncs []func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)
}

// GatewayConfig 网关服务配置
type GatewayConfig struct {
	GateAddr          string
	UseProtoNames     bool
	AuthPublicKeyFile string
	GrpcSubConfigs    []GrpcGatewayConfig
	HandlePathConfigs []HandlePathConfig
}

type HandlePathConfig struct {
	Meth        string
	PathPattern string
	Handle      runtime.HandlerFunc
}

func customHeaderMatcher(key string) (string, bool) {
	if strings.HasPrefix(key, "X-") {
		return key, true
	}
	return runtime.DefaultHeaderMatcher(key)
}

// RunGatewayServer 启动网关服务
func RunGatewayServer(config *GatewayConfig) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(customHeaderMatcher),
		runtime.WithMarshalerOption(
			runtime.MIMEWildcard,
			&runtime.JSONPb{
				MarshalOptions: protojson.MarshalOptions{
					UseEnumNumbers:  true,
					UseProtoNames:   config.UseProtoNames,
					EmitUnpopulated: true,
				},
				UnmarshalOptions: protojson.UnmarshalOptions{
					DiscardUnknown: true, // If DiscardUnknown is set, unknown fields are ignored.
				},
			},
		))
	if len(config.AuthPublicKeyFile) > 0 {

	}
	for _, h := range config.HandlePathConfigs {
		err := mux.HandlePath(h.Meth, h.PathPattern, h.Handle)
		if err != nil {
			panic("gateway cannot HandlePath: " + err.Error())
		}
	}
	for _, c := range config.GrpcSubConfigs {
		for _, f := range c.RegisterFuncs {
			err := f(ctx, mux, c.ServerAddr,
				[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
			)
			if err != nil {
				panic("gateway cannot register service: " + err.Error())
			}
		}
	}
	log.Infof("grpc gateway started at %s", config.GateAddr)
	panic(http.ListenAndServe(config.GateAddr, tracingWrapper(mux)))
}

func tracingWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debugf("ip: %v", exnet.ClientPublicIP(r))
		body, err := httputil.DumpRequest(r, true)
		if err != nil {
			log.Errorf("httputil.DumpRequest err: %v", err)
		} else {
			if !strings.HasPrefix(r.URL.Path, "/upload") {
				log.Debugf("[body]%s", body)
			}
		}
		h.ServeHTTP(w, r)
	})
}
