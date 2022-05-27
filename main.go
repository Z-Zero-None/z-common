package main

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"z-common/global"
	"z-common/tools/middleware"

	_ "z-common/setup"
)

func DoSomething(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"done": "今天要努力!!!,zhong ze nan",
	})
}

func main() {
	r := gin.Default()
	r.Use(middleware.JaegerTracing(global.JaegerTrace))
	r.Use(middleware.PrometheusMonitoring())
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.GET("/done", DoSomething)
	err := r.Run(":8888")
	if err != nil {
		return
	}
}
