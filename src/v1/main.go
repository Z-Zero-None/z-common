package main

import (
	database "cloud.google.com/go/spanner/admin/database/apiv1"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	_ "z-common/src/base/setup"
	"z-common/src/v1/base/global"
	"z-common/src/v1/base/middleware"
)

type Country struct {
	Code       string
	Name       string
	Population int
	Age        int
}

func DoSomething(ctx *gin.Context) {
	engine, err := database.NewDefaultMysqlConfig().GetMySQLEngine()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"errMsg": err.Error(),
		})
	}
	country := Country{}
	if err = engine.Table("country").Limit(1).Scan(&country).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"errMsg": err.Error(),
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"done":   "今天要努力!!!,zhong ze nan",
		"dbData": fmt.Sprintf("country:%v", country),
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
