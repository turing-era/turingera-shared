package log

import (
	"context"
	"runtime"
	"time"

	"google.golang.org/grpc"

	"github.com/turing-era/turingera-shared/auth"
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

	userID, _ := auth.UserIDFromCtx(ctx)
	if err != nil {
		Errorf("[serverlog][userid:%s]: %s, cost: %v, req: %s, err: %v", userID, info.FullMethod, cost, reqJson, err)
	} else {
		Debugf("[serverlog][userid:%s]: %s, cost: %v, req: %s, rsp: %s", userID, info.FullMethod, cost, reqJson, cutils.Obj2Json(rsp))
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
