package server

import (
	"context"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/turing-era/turingera-shared/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

// GatewaySubConfig 网关服务子配置
type GatewaySubConfig struct {
	ServerAddr    string
	RegisterFuncs []func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)
}

// GatewayConfig 网关服务配置
type GatewayConfig struct {
	GateAddr      string
	UseProtoNames bool
	SubConfigs    []GatewaySubConfig
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
	for _, c := range config.SubConfigs {
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
		body, _ := httputil.DumpRequest(r, true)
		log.Debugf("[%v] %s", r.URL, body)
		h.ServeHTTP(w, r)
	})
}
