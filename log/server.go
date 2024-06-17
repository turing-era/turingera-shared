package log

import (
	"context"
	"runtime"
	"time"

	"google.golang.org/grpc"

	"github.com/turing-era/turingera-shared/cutils"
)

const panicBufLen = 1024

// ServerLogInterceptor 拦截器方法
func ServerLogInterceptor(ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	defer ServerRecover()

	reqJson := cutils.Obj2Json(req)
	rsp, err := handler(ctx, req)
	cost := time.Since(start)
	if err != nil {
		Errorf("%s, cost: %v, req: %s, err: %v", info.FullMethod, cost, reqJson, err)
	} else {
		Debugf("%s, cost: %v, req: %s, rsp: %s", info.FullMethod, cost, reqJson, cutils.Obj2Json(rsp))
	}
	return rsp, err
}

func ServerRecover() {
	if e := recover(); e != nil {
		buf := make([]byte, panicBufLen)
		buf = buf[:runtime.Stack(buf, false)]
		Errorf("[PANIC]%v\n%s\n", e, buf)
	}
}
