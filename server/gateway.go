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

// GatewayConfig 网关服务配置
type GatewayConfig struct {
	Addr         string
	RegisterFunc func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)
}

func customHeaderMatcher(key string) (string, bool) {
	if strings.HasPrefix(key, "X-") {
		return key, true
	}
	return runtime.DefaultHeaderMatcher(key)
}

// RunGatewayServer 启动网关服务
func RunGatewayServer(configs []*GatewayConfig) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(customHeaderMatcher),
		runtime.WithMarshalerOption(
			runtime.MIMEWildcard,
			&runtime.JSONPb{
				MarshalOptions: protojson.MarshalOptions{
					UseEnumNumbers: true,
					// UseProtoNames:   true,
					EmitUnpopulated: true,
				},
				UnmarshalOptions: protojson.UnmarshalOptions{
					DiscardUnknown: true, // If DiscardUnknown is set, unknown fields are ignored.
				},
			},
		))
	for _, c := range configs {
		err := c.RegisterFunc(ctx, mux, c.Addr,
			[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
		)
		if err != nil {
			panic("gateway cannot register service: " + err.Error())
		}
	}
	addr := ":8081"
	log.Infof("grpc gateway started at %s", addr)
	panic(http.ListenAndServe(addr, tracingWrapper(mux)))
}

func tracingWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := httputil.DumpRequest(r, true)
		log.Debugf("[%v] %s", r.URL, body)
		h.ServeHTTP(w, r)
	})
}
