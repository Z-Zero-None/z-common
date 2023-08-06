FROM golang:1.20 as builder

# 设置工作目录
WORKDIR /go/src/app

# 复制文件
COPY . .

ENV GOPROXY https://goproxy.cn,direct
# $GOPATH/bin添加到环境变量中
ENV PATH $GOPATH/bin:$PATH
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
RUN go mod tidy
RUN go build -ldflags="-s -w" -o app main.go

FROM alpine:3.14.0
WORKDIR /go/src/app
# 把执行builder阶段的结果 /app/main拷贝到/app中
COPY --from=builder /go/src/app/app .
COPY docker/app/run_app.sh .
RUN chmod +x ./run_app.sh

# 暴露端口
EXPOSE 8080
ENTRYPOINT ["/bin/sh"]
CMD ["run_app.sh"]