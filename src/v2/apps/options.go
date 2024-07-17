package apps

import (
	"bytes"
	"crypto/md5"
	"encoding/gob"
	"fmt"
)

const (
	DebugMode = iota + 1
	TestMode
	ReleaseMode
)

var mode = DebugMode

func SetMode(env string) {
	switch env {
	case "dev":
		mode = DebugMode
	case "test":
		mode = TestMode
	default:
		mode = ReleaseMode
	}
}

func Mode() int {
	return mode
}

var namespace = ""

func Sprintf(format string, v ...interface{}) string {
	return fmt.Sprintf(namespace+format, v...)
}

func SprintHash(prefix string, v interface{}) string {
	buf := bytes.NewBuffer(nil)
	err := gob.NewEncoder(buf).Encode(v)
	var keyStr string
	if err == nil {
		h := md5.Sum(buf.Bytes())
		keyStr = Sprintf(prefix+"%x", h[6:])
	}
	return keyStr
}
