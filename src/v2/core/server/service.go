package server

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

type Service struct {
	Options *Options
}

func New(opts ...Option) *Service {
	sopts := Options{
		Singles: []os.Signal{syscall.SIGKILL, syscall.SIGINT, syscall.SIGABRT, syscall.SIGTERM},
		Servers: nil,
	}

	for _, o := range opts {
		o(&sopts)
	}

	return &Service{
		Options: &sopts,
	}
}

func (s *Service) Run() error {
	if len(s.Options.Servers) == 0 {
		return nil
	}

	errCh := make(chan error)
	for _, srv := range s.Options.Servers {
		go func(s Server) {
			errCh <- s.Run()
		}(srv)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, s.Options.Singles...)

	select {
	case err := <-errCh:
		if err != nil {
			slog.Error("service.Run", "error", err)
		}
		s.Close()
		return err
	case ch := <-sig:
		slog.Info("service.Run", "signal", ch)
		s.Close()
		slog.Info("service.Run Exit!")
	}

	return nil
}

func (s *Service) Close() error {
	for _, srv := range s.Options.Servers {
		srv.Close()
	}
	for _, fn := range s.Options.CloseFuncs {
		fn()
	}

	return nil
}
