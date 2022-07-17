# CloudNative

# 第三次作业

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

## 2、创建 Service 通过 nodePort 提供对外访问的端口

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

## 3、创建 configmap 提供环境变量和配置等

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

## 4、创建服务实例

- 优雅启动

- 优雅终止

- 资源需求和 QoS 保证

- 探活

- 日常运维需求，日志等级

- 配置和代码分离
