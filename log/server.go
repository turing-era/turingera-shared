package log

import (
	"context"
	"time"

	"github.com/turing-era/turingera/cutils"
	"google.golang.org/grpc"
)

// ServerLogInterceptor 拦截器方法
func ServerLogInterceptor(ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	reqJson := cutils.Obj2Json(req)
	start := time.Now()
	rsp, err := handler(ctx, req)
	rspJson := cutils.Obj2Json(rsp)
	cost := time.Since(start).Milliseconds()
	if err == nil {
		Debugf("[%v][cost: %vms] req: %v, rsp: %v", info.FullMethod, cost, reqJson, rspJson)
	} else {
		Errorf("[%v][cost: %vms] req: %v, err: %v", info.FullMethod, cost, reqJson, err)
	}
	return rsp, err
}
