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