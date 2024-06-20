package core

import (
	"context"
	"os"
	"time"
	"z-common/src/v2/core/transport"
)

type Option func(o *options)
type CloseFunc func()

type options struct {
	ctx         context.Context
	id          string        //唯一标识
	name        string        //服务模块名称
	singles     []os.Signal   //监听信号处理相关问题
	stopTimeout time.Duration //停止服务时间处理
	servers     []transport.Server
	beforeStart []func(context.Context) error
	afterStart  []func(context.Context) error
	beforeStop  []func(context.Context) error
	afterStop   []func(context.Context) error
	closeFns    []CloseFunc
}

func WithContext(ctx context.Context) Option {
	return func(o *options) { o.ctx = ctx }
}

func WithID(id string) Option {
	return func(o *options) { o.id = id }
}

func WithName(name string) Option {
	return func(o *options) { o.name = name }
}

func WithSignals(sigs ...os.Signal) Option {
	return func(o *options) {
		o.singles = append(o.singles, sigs...)
	}
}

func WithStopTimeout(t time.Duration) Option {
	return func(o *options) { o.stopTimeout = t }
}

func WithServer(s ...transport.Server) Option {
	return func(o *options) {
		o.servers = append(o.servers, s...)
	}
}

func WithBeforeStart(fn func(context.Context) error) Option {
	return func(o *options) {
		o.beforeStart = append(o.beforeStart, fn)
	}
}

func WithBeforeStop(fn func(context.Context) error) Option {
	return func(o *options) {
		o.beforeStop = append(o.beforeStop, fn)
	}
}

func WithAfterStart(fn func(context.Context) error) Option {
	return func(o *options) {
		o.afterStart = append(o.afterStart, fn)
	}
}

func WithAfterStop(fn func(context.Context) error) Option {
	return func(o *options) {
		o.afterStop = append(o.afterStop, fn)
	}
}

func WithCloseFns(fns ...CloseFunc) Option {
	return func(o *options) {
		o.closeFns = append(o.closeFns, fns...)
	}
}
