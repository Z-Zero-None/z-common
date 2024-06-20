package core

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type appKey struct{}

type AppInfo interface {
	Empty() string
}

type App struct {
	opts   options
	ctx    context.Context
	cancel context.CancelFunc
}

func (a *App) Empty() string {
	return "app"
}

func New(opts ...Option) *App {
	o := options{
		ctx:         context.Background(),
		singles:     []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGKILL, syscall.SIGABRT},
		stopTimeout: 10 * time.Second,
	}
	if id, err := uuid.NewUUID(); err == nil {
		o.id = id.String()
	}
	for _, opt := range opts {
		opt(&o)
	}
	ctx, cancel := context.WithCancel(o.ctx)
	return &App{
		ctx:    ctx,
		cancel: cancel,
		opts:   o,
	}
}

func NewContext(ctx context.Context, s AppInfo) context.Context {
	return context.WithValue(ctx, appKey{}, s)
}

func (a *App) Run() error {
	if len(a.opts.servers) == 0 {
		slog.Warn("App.Run: empty servers")
		return nil
	}
	sctx := NewContext(a.ctx, a)
	eg, ctx := errgroup.WithContext(sctx)
	wg := sync.WaitGroup{}
	for _, fn := range a.opts.beforeStart {
		if err := fn(sctx); err != nil {
			return err
		}
	}
	for _, srv := range a.opts.servers {
		srv := srv
		eg.Go(func() error {
			<-ctx.Done() // wait for stop signal
			stopCtx, cancel := context.WithTimeout(NewContext(a.opts.ctx, a), a.opts.stopTimeout)
			defer cancel()
			return srv.Stop(stopCtx)
		})
		wg.Add(1)
		eg.Go(func() error {
			wg.Done() // here is to ensure server start has begun running before register, so defer is not needed
			return srv.Start(NewContext(a.opts.ctx, a))
		})
	}
	wg.Wait()
	//todo 注册服务

	for _, fn := range a.opts.afterStart {
		if err := fn(sctx); err != nil {
			return err
		}
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, a.opts.singles...)
	eg.Go(func() error {
		select {
		case <-ctx.Done(): //取消触发
			slog.Error("service.Run.cancel")
			return nil
		case ch := <-c: //信号触发
			slog.Error("service.Run", "signal", ch)
			return a.Stop()
		}
	})
	if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	var err error
	for _, fn := range a.opts.afterStop {
		err = fn(sctx)
	}
	return err
}

func (a *App) Stop() (err error) {
	sctx := NewContext(a.ctx, a)
	for _, fn := range a.opts.beforeStop {
		err = fn(sctx)
	}
	//关闭在外所有服务
	if a.cancel != nil {
		a.cancel()
	}
	for _, fn := range a.opts.closeFns {
		fn()
	}
	return err
}
