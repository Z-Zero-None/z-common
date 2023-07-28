package serve

import "github.com/gin-gonic/gin"

type IHttpServer interface {
	AddMiddleware(list []func() gin.HandlerFunc)
	AddHandler(info *HandlerInfo) error
	Run(port int) error
}
