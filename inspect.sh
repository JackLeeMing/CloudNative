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
