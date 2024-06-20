package httpgin

import (
	"context"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net"
	"net/http"
)

type ServerOption func(*Server)

type Server struct {
	engine     *gin.Engine
	addr       string
	name       string
	httpServer *http.Server
	//todo 服务发现 socket处理
}

func WithName(name string) ServerOption {
	return func(s *Server) {
		s.name = name
	}
}

func WithAddr(addr string) ServerOption {
	return func(s *Server) {
		s.addr = addr
	}
}

func NewServer(options ...ServerOption) *Server {
	engine := gin.Default()
	// 设置默认起始中间件
	//engine.Use(func(c *gin.Context) {
	//	c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "trace_id", uuid.New().String()))
	//	c.Next()
	//})
	svc := &Server{
		engine:     engine,
		httpServer: &http.Server{Handler: engine},
	}
	for _, o := range options {
		o(svc)
	}
	return svc
}

func (s *Server) Engine() *gin.Engine {
	return s.engine
}

func (s *Server) Start(ctx context.Context) error {
	ln, err := net.Listen("tcp", s.addr)
	s.httpServer.BaseContext = func(ln net.Listener) context.Context {
		return ctx
	}
	if err != nil {
		return err
	}
	slog.Info("[HTTP] server listening on: %s", ln.Addr().String())
	return s.httpServer.Serve(ln)
}

func (s *Server) Stop(ctx context.Context) error {
	slog.Info("[HTTP] server stopping")
	return s.httpServer.Shutdown(ctx)
}
