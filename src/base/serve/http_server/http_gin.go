package http_server

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
	"z-common/src/base/serve"
)

var middleWareList = []gin.HandlerFunc{
	corsMiddleware(), recordBaseInfo(),
}

type httpGinServer struct {
	engine *gin.Engine
}

func (h *httpGinServer) addMiddleware(list []gin.HandlerFunc) {
	for _, f := range list {
		h.engine.Use(f)
	}
}

func NewServer() serve.IHttpServer {
	app := gin.Default()
	server := &httpGinServer{
		engine: app,
	}
	server.addMiddleware(middleWareList)
	return server
}

func (h *httpGinServer) AddMiddleware(list []gin.HandlerFunc) {
	h.addMiddleware(list)
}

func (h *httpGinServer) AddHandler(info *serve.HandlerInfo) error {
	params := make([]gin.HandlerFunc, 0)
	for _, handler := range info.MiddlewareHandlers {
		params = append(params, handler.(func(ctx *gin.Context)))
	}
	params = append(params, info.Handler.(func(ctx *gin.Context)))
	switch info.Method {
	case http.MethodGet:
		h.engine.GET(info.Path, params...)
	case http.MethodPost:
		h.engine.POST(info.Path, params...)
	case http.MethodPut:
		h.engine.PUT(info.Path, params...)
	case http.MethodDelete:
		h.engine.DELETE(info.Path, params...)
	case http.MethodPatch:
		h.engine.PATCH(info.Path, params...)
	default:
		return errors.New("Invalid method of HandlerInfo")
	}
	return nil
}

func (h *httpGinServer) Run(port string) error {
	readTimeout, _ := strconv.Atoi(os.Getenv("serve_read_timeout"))
	writeTimeout, _ := strconv.Atoi(os.Getenv("serve_write_timeout"))
	s := &http.Server{
		Addr:           ":" + port,
		Handler:        h.engine,
		ReadTimeout:    time.Duration(readTimeout) * time.Second,
		WriteTimeout:   time.Duration(writeTimeout) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	//优雅重启
	go func() {
		err := s.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("s.ListenAndServe err:%v", err)
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	timeout, _ := strconv.Atoi(os.Getenv("serve_shutdown_timeout"))
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("Server force to shutdown err:%v", err)
	}
	log.Println("Server exiting!!!")
	return nil
}
