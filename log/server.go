package log

import (
	"context"
	"runtime"
	"time"

	"github.com/turing-era/turingera-shared/cutils"
	"google.golang.org/grpc"
)

const panicBufLen = 1024

// ServerLogInterceptor 拦截器方法
func ServerLogInterceptor(ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	defer func() {
		if e := recover(); e != nil {
			buf := make([]byte, panicBufLen)
			buf = buf[:runtime.Stack(buf, false)]
			Errorf("[PANIC]%v\n%s\n", e, buf)
		}
	}()
	reqJson := cutils.Obj2Json(req)
	Debugf("[%v]req: %v", info.FullMethod, reqJson)
	start := time.Now()

	rsp, err := handler(ctx, req)
	cost := time.Since(start).Milliseconds()
	if err != nil {
		Errorf("[%v][cost: %vms]err: %v", info.FullMethod, cost, err)
	} else {
		rspJson := cutils.Obj2Json(rsp)
		Debugf("[%v][cost: %vms]rsp: %v", info.FullMethod, cost, rspJson)
	}
	return rsp, err
}
