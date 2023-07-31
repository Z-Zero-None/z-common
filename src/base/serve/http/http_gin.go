package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
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

func (h *httpGinServer) Run(port int) error {
	addr := fmt.Sprintf(":%d", port)
	err := h.engine.Run(addr)
	if err != nil {
		return errors.Wrap(err, "Failed to run http server")
	}
	return nil
}
