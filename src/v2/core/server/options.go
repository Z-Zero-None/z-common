package server

import (
	"os"
)

type Option func(options *Options)
type CloseFunc func()

type Options struct {
	Singles    []os.Signal
	Servers    []Server
	CloseFuncs []CloseFunc
}

func WithServer(s ...Server) Option {
	return func(opts *Options) {
		opts.Servers = append(opts.Servers, s...)
	}
}

func WithCloseFunc(funcs ...CloseFunc) Option {
	return func(opts *Options) {
		opts.CloseFuncs = append(opts.CloseFuncs, funcs...)
	}
}

func WithSignal(sig ...os.Signal) Option {
	return func(opts *Options) {
		opts.Singles = append(opts.Singles, sig...)
	}
}
