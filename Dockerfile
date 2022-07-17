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

FROM jackleeming/cloudnative:v1-tini

COPY --from=builder /go/bin/CloudNative /go/bin/CloudNative

WORKDIR /go/bin

EXPOSE 8090

CMD ["./CloudNative"]