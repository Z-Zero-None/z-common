package http_server

import (
	"bytes"
	"github.com/apolloconfig/agollo/v4/component/log"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strconv"
	"time"
)

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if gin.Mode() == gin.DebugMode && c.Request.Method == "OPTIONS" {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
			c.Header("Access-Control-Allow-Methods", "GET, OPTIONS, POST, PUT, DELETE")
			c.Set("Content-Type", "application/json")
			c.JSON(http.StatusOK, "")
			return
		}
		c.Next()
	}
}

func recordBaseInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		bodyLen, _ := strconv.Atoi(c.Request.Header.Get("content-length"))
		var bodyBytes []byte
		// avoid to log too large request body.
		if c.Request.Method == http.MethodPost && bodyLen < 4096 {
			// Read the Body content
			if c.Request.Body != nil {
				bodyBytes, _ = io.ReadAll(c.Request.Body)
			}
			// Restore the io.ReadCloser to its original state
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}
		requestParams := string(bodyBytes)
		start := time.Now()
		c.Next()
		since := time.Since(start)
		record := map[string]interface{}{
			"host":      c.Request.Host,
			"client_ip": c.Request.RemoteAddr,
			"method":    c.Request.Method,
			"path":      c.Request.RequestURI,
			"request":   requestParams,
			"since":     since.String(),
		}
		log.Info(record)
	}
}
