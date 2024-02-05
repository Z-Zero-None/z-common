package setup

import (
	"fmt"
	"z-common/global"
)

func init() {
	setupJaegerTrace()
}

func setupJaegerTrace() error {
	config := database.NewDefaultJaegerTraceConfig()
	jaegerTrace, _, err := database.NewJaegerTrace(config)
	if err != nil {
		return fmt.Errorf("setupJaegerTrace err:%v", err)
	}
	global.JaegerTrace = jaegerTrace
	return err
}
