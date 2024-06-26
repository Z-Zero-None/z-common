package main

import (
	"z-common/src/v1/base/serve/http_server"
)

func main() {
	server := http_server.NewServer()
	server.Run("8888")
}
