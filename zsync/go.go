package zsync

import (
	"context"
	"runtime"
	"time"

	"github.com/panjf2000/ants/v2"
	"github.com/turing-era/turingera-shared/log"
	"google.golang.org/grpc/metadata"
)

// Goer is the interface that launches a testable and safe goroutine.
type Goer interface {
	Go(ctx context.Context, timeout time.Duration, handler func(context.Context)) error
}

type asyncGoer struct {
	panicBufLen   int
	shouldRecover bool
	pool          *ants.PoolWithFunc
}

type goerParam struct {
	ctx     context.Context
	cancel  context.CancelFunc
	handler func(context.Context)
}

// NewAsyncGoer creates a goer that executes handler asynchronously with a goroutine when Go() is called.
func NewAsyncGoer(workerPoolSize int, panicBufLen int, shouldRecover bool) Goer {
	g := &asyncGoer{
		panicBufLen:   panicBufLen,
		shouldRecover: shouldRecover,
	}
	if workerPoolSize == 0 {
		return g
	}

	pool, err := ants.NewPoolWithFunc(workerPoolSize, func(args interface{}) {
		p := args.(*goerParam)
		g.handle(p.ctx, p.handler, p.cancel)
	})
	if err != nil {
		panic(err)
	}
	g.pool = pool
	return g
}

func (g *asyncGoer) handle(ctx context.Context, handler func(context.Context), cancel context.CancelFunc) {
	defer func() {
		if g.shouldRecover {
			if err := recover(); err != nil {
				buf := make([]byte, g.panicBufLen)
				buf = buf[:runtime.Stack(buf, false)]
				log.Errorf("[PANIC]%v\n%s\n", err, buf)
			}
		}
		cancel()
	}()
	handler(ctx)
}

func (g *asyncGoer) Go(ctx context.Context, timeout time.Duration, handler func(context.Context)) error {
	md, _ := metadata.FromOutgoingContext(ctx)
	newCtx := context.Background()
	newCtx = metadata.NewOutgoingContext(newCtx, md)
	newCtx, cancel := context.WithTimeout(newCtx, timeout)
	if g.pool != nil {
		p := &goerParam{
			ctx:     newCtx,
			cancel:  cancel,
			handler: handler,
		}
		return g.pool.Invoke(p)
	}
	go g.handle(newCtx, handler, cancel)
	return nil
}

// DefaultGoer is an async goer without workerpool.
var DefaultGoer = NewAsyncGoer(0, PanicBufLen, true)

// Go launches a safer goroutine for async task inside rpc handler.
// it clones ctx and msg before the goroutine, and will recover and report metrics when the goroutine panics.
// you should set a suitable timeout to control the lifetime of the new goroutine to prevent goroutine leaks.
func Go(ctx context.Context, timeout time.Duration, handler func(context.Context)) error {
	return DefaultGoer.Go(ctx, timeout, handler)
}
