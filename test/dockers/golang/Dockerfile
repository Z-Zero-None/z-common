FROM golang:1.17

ENV GOPROXY https://goproxy.cn,direct
# $GOPATH/bin添加到环境变量中
ENV PATH $GOPATH/bin:$PATH
# 设置工作目录
WORKDIR /usr/src/code
# 更新驱动
USER root
RUN apt-get update && apt-get install -y vim
# 清理垃圾
RUN apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* && \
    rm /var/log/lastlog /var/log/faillog
# 暴露端口
EXPOSE 8888