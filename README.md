# CloudNative

# 第三次作业

# 第一部分

## 1、 创建 namespace

```Yaml
apiVersion: v1
kind: Namespace
metadata:
  labels:
    k8s-app: cloudnative
    cloudnative/name: cloudnative
  name: cloudnative
```

## 2、创建 configmap 提供环境变量和配置等

```Yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: config-env
data:
  version: v1.0.3
  loglevel: debug
  httpport: "8090"
```

## 3、创建服务实例

```Yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: httpserver
  namespace: cloudnative
spec:
  replicas: 1
  selector:
    matchLabels:
      k8s-app: cloudnative
      cloudnative/name: cloudnative
  template:
    metadata:
      labels:
        k8s-app: cloudnative
        cloudnative/name: cloudnative
    spec:
      containers:
      - image: jackleeming/cloudnative:v1.0.4
        imagePullPolicy: IfNotPresent
        # 加载配置
        envFrom:
        - configMapRef:
            name: config-env
        #  liveness 探针
        livenessProbe:
            failureThreshold: 3
            httpGet:
              path: /healthz
              port: 8090
              scheme: HTTP
            initialDelaySeconds: 5
            periodSeconds: 10
            successThreshold: 1
            timeoutSeconds: 1
        name: httpserver
        #  readiness 探针
        readinessProbe:
            failureThreshold: 3
            httpGet:
              path: /healthz
              port: 8090
              scheme: HTTP
            initialDelaySeconds: 5
            periodSeconds: 10
            successThreshold: 1
            timeoutSeconds: 1
        # 资源配置
        resources:
          limits:
            cpu: 500m
            memory: 500Mi
          requests:
            cpu: 200m
            memory: 200Mi
        volumeMounts:
        - mountPath: /etc/localtime
          name: localtime
        ports:
        - containerPort: 8090
      restartPolicy: Always
      terminationGracePeriodSeconds: 30
      volumes:
      - name: localtime
        hostPath:
          path: /etc/localtime
```

## 4、关于容器终止

- 1、构建基于 tini 的 alpine 基础镜像 jackleeming/cloudnative:v1-tini

```Dockerfile
FROM alpine:3.16.0

RUN apk add --no-cache tini
# Tini is now available at /sbin/tini
ENTRYPOINT ["/sbin/tini", "--"]
```

- 2、修改 Dockerfile，第二阶段的引用镜像设置为自定义的 jackleeming/cloudnative:v1-tini 镜像，用 tini 来管理容器中的进程

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
# 修改第二阶段的基础镜像
FROM jackleeming/cloudnative:v1-tini

COPY --from=builder /go/bin/CloudNative /go/bin/CloudNative

WORKDIR /go/bin

EXPOSE 8090
# 程序启动采用 CMD，通过tini来启动
CMD ["./CloudNative"]
```

## 完整的服务 部署文件见 deployment.yml

## 第二部分

## 基于 Service 通过 nodePort 提供对外访问的端口

```Yaml
apiVersion: v1
kind: Service
metadata:
  name: httpserver
  namespace: cloudnative
  labels:
    k8s-app: cloudnative
    cloudnative/name: cloudnative
spec:
  ports:
  - name: httpserver
    port: 8090
    protocol: TCP
    targetPort: 8090
    nodePort: 32001
  selector:
    k8s-app: cloudnative
    cloudnative/name: cloudnative
  type: NodePort
```

## 基于 ingress-nignx

- 1、ingress 部署文件见 ./ingress/deploy_1.3.0.yaml 镜像采用 registry 搭建的私有非加密仓库
- 2、创建 证书

```sh
openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout tls.key -out tls.crt -subj "/CN=cncamp.com/O=cncamp" -addext "subjectAltName = DNS:cncamp.com"
```

- 3、创建 secret

```sh
kubectl create secret tls cncamp-tls --cert=./tls.crt --key=./tls.key
```

- 4、创建 ingress

```Yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: httpserver
  namespace: cloudnative
  annotations:
    kubernetes.io/ingress.class: "nginx"
spec:
  tls:
    - hosts:
        - cncamp.com
      secretName: cncamp-tls
  rules:
    - host: cncamp.com
      http:
        paths:
          - path: "/"
            pathType: Prefix
            backend:
              service:
                name: httpserver
                port:
                  number: 8090
```

- 5 结果演示

  ![](./p1.png)

  ![](./p2.png)

  ![](./p3.png)

## 完整的服务部署文件见 deployment.yml
