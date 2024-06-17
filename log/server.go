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

	rsp, err := handler(ctx, req)
	cost := time.Since(start).Milliseconds()
	if err != nil {
		Errorf("[%v][cost: %vms]req: %v\nerr: %v", info.FullMethod, cost, cutils.Obj2Json(req), err)
	} else {
		Debugf("[%v][cost: %vms]req: %v\nrsp: %v", info.FullMethod, cost, cutils.Obj2Json(req), cutils.Obj2Json(rsp))
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
