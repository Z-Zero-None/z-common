package http

import (
	"github.com/gin-gonic/gin"
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
}

func (h *httpGinServer) Run(port int) error {
	//TODO implement me
	panic("implement me")
}
