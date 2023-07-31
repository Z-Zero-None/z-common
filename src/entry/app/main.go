package main

import "z-common/src/base/serve/http_server"

func main() {
	server := http_server.NewServer()
	server.Run("8888")
}
