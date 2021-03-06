# CloudNative

# 第二次作业

## 镜像构建

- 基础镜像

```Dockerfile
FROM golang:1.17.11-alpine3.16

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories \
    && apk update \
    && apk add --no-cache git  \
    && apk add --no-cache docker    \
    && apk add --no-cache tzdata    \
    && apk add --no-cache curl  \
    && apk add --no-cache gettext	\
    && apk add --no-cache cloc  \
    && apk add --no-cache upx

RUN go env -w GO111MODULE=on \
    && go env -w CGO_ENABLED=0 \
    && go env -w GOPROXY=https://goproxy.cn,direct

ENV TZ=Asia/Shanghai
```

- 工程镜像

```Dockerfile
FROM jackleeming/cloudnative:v1-base AS builder

RUN mkdir -p $GOPATH/src/github.com/CloudNative

WORKDIR $GOPATH/src/github.com/CloudNative

COPY . .

RUN go env -w GO111MODULE=on \
    && go env -w CGO_ENABLED=0 \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && go mod download \
    && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags '-w -s -extldflags "-static"' -o /go/bin/CloudNative main.go

WORKDIR /go/bin

RUN upx CloudNative

FROM alpine:3.16.0

COPY --from=builder /go/bin/CloudNative /go/bin/CloudNative

WORKDIR /go/bin

EXPOSE 8090

ENTRYPOINT ["./CloudNative"]
```

## nsenter 操作步骤

```shell
echo "#1 启动容器"
docker run -d --name pss1 -p 8090:8090 jackleeming/cloudnative:v1.0.1
echo "# 获取容器中的进程id"
docker inspect pss1 | grep -i pid
################################
pid=$(docker inspect --format "{{.State.Pid}}" pss1)
echo "# 查看容器的路由"
nsenter -t $pid -n ip r
echo "# 查看容器的 addr"
nsenter -t $pid -n ip addr
```
