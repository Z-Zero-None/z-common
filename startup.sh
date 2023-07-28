#!/bin/bash
#编译文件为main可执行文件
GOOS="linux" CGO_ENABLED="0" go build -ldflags="-s -w"-o main main.go
#执行main服务
./main