package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

func init() {
	prometheus.MustRegister(httpRequestCounter)
	prometheus.MustRegister(httpRequestsHistogram)
}

var httpRequestCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "zzn_request_count",
		Help: "request count",
	},
)

// 监控请求量，请求耗时等
var httpRequestsHistogram = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Subsystem: "sdk",
		Name:      "zzn_handle_requests",
		Help:      "Histogram statistics of http requests handle by elete http. Buckets by latency",
		Buckets:   []float64{0.001, 0.002, 0.005, 0.01, 0.05, 0.1, 0.2, 0.3, 0.4, 0.5, 0.8, 1, 2, 5, 10},
	},
	[]string{"code"},
)

func PrometheusMonitoring() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 统计接口请求数量
		httpRequestCounter.Inc()
		startTime := time.Now()
		// 处理后续逻辑
		c.Next()
		// after request
		finishTime := time.Now()
		// 监控计算接口耗时，请求数量等
		httpRequestsHistogram.With(prometheus.Labels{"code": strconv.Itoa(c.Writer.Status())}).Observe(float64(finishTime.Sub(startTime)) / (1000 * 1000 * 1000))
	}
}
