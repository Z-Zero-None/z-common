package serve

import "github.com/gin-gonic/gin"

type IHttpServer interface {
	AddMiddleware(list []gin.HandlerFunc)
	AddHandler(info *HandlerInfo) error
	Run(port string) error
}
