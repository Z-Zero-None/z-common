package setup

import (
	"fmt"
	"z-common/global"
	"z-common/src/base/connector"
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
