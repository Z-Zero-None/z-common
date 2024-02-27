package setup

import (
	"fmt"
	"z-common/global"
	"z-common/src/base/connectors"
)

func init() {
	setupJaegerTrace()
}

func setupJaegerTrace() error {
	config := connectors.NewDefaultJaegerTraceConfig()
	jaegerTrace, _, err := connectors.NewJaegerTrace(config)
	if err != nil {
		return fmt.Errorf("setupJaegerTrace err:%v", err)
	}
	global.JaegerTrace = jaegerTrace
	return err
}
