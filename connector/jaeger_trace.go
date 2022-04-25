package connector

import (
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"io"
	"time"
)

type JaegerTraceConfig struct {
	ServiceName            string
	SamplerType            string
	SamplerParam           float64
	ReporterLogSpans       bool
	ReporterLocalAgentHost string
}

var defaultJaegerTraceConfig = JaegerTraceConfig{
	ServiceName:            "test",
	SamplerType:            jaeger.SamplerTypeConst,
	SamplerParam:           1,
	ReporterLogSpans:       true,
	ReporterLocalAgentHost: "127.0.0.1:6831",
}

func NewDefaultJaegerTraceConfig() *JaegerTraceConfig {
	return &defaultJaegerTraceConfig
}

func NewJaegerTrace(setting *JaegerTraceConfig) (opentracing.Tracer, io.Closer, error) {
	cfg := config.Configuration{
		ServiceName: setting.ServiceName,
		Sampler: &config.SamplerConfig{
			Type:  setting.SamplerType,
			Param: setting.SamplerParam,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            setting.ReporterLogSpans,
			LocalAgentHostPort:  setting.ReporterLocalAgentHost,
			BufferFlushInterval: time.Second * 1,
		},
	}
	tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
	if err != nil {
		return nil, nil, err
	}
	return tracer, closer, err
}
