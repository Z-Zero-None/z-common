package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func DoSomething(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"done": "今天要努力!!!",
	})
}

func main() {
	r := gin.Default()
	r.GET("/done", DoSomething)
	err := r.Run(":8888")
	if err != nil {
		return
	}
}
