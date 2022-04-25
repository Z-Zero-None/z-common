package setup

import (
	"fmt"
	"z-common/connector"
	"z-common/global"
)

func init() {
	setupJaegerTrace()
}

func setupJaegerTrace() error {
	config := connector.NewDefaultJaegerTraceConfig()
	jaegerTrace, _, err := connector.NewJaegerTrace(config)
	if err != nil {
		return fmt.Errorf("setupJaegerTrace err:%v", err)
	}
	global.JaegerTrace = jaegerTrace
	return err
}
