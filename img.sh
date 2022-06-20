docker tag k8s.gcr.io/kube-apiserver:v1.24.2 docker-registry.jaquelee.com:5000/kube-apiserver:v1.24.2
docker tag k8s.gcr.io/kube-proxy:v1.24.2 docker-registry.jaquelee.com:5000/kube-proxy:v1.24.2
docker tag k8s.gcr.io/kube-scheduler:v1.24.2 docker-registry.jaquelee.com:5000/kube-scheduler:v1.24.2
docker tag k8s.gcr.io/kube-controller-manager:v1.24.2 docker-registry.jaquelee.com:5000/kube-controller-manager:v1.24.2

docker tag k8s.gcr.io/etcd:3.5.3-0 docker-registry.jaquelee.com:5000/etcd:3.5.3-0
docker tag k8s.gcr.io/pause:3.7 docker-registry.jaquelee.com:5000/pause:3.7
docker tag k8s.gcr.io/coredns/coredns:v1.8.6 docker-registry.jaquelee.com:5000/coredns/coredns:v1.8.6
