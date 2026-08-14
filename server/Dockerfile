FROM golang:1.23.2-alpine AS builder

WORKDIR /server

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w -s" -o apiserver cmd/apiserver/run.go

# 声明镜像来源为alpine:latest
FROM alpine:latest
LABEL MAINTAINER="ivan.xiewei@foxmail.com"
# 设置时区
ENV TZ=Asia/Shanghai
# 声明工作目录
WORKDIR /server

# 把可执行文件和配置文件复制到src
COPY --from=builder /server/apiserver ./
COPY --from=builder /server/config.release.yaml ./
RUN mkdir db

EXPOSE 8888

# 运行打包好的二进制 并用-c 指定config.release.yaml配置文件
ENTRYPOINT ["./apiserver", "-c",  "config.release.yaml"]