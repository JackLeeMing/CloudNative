curl -H "Host: cloudnative.jaquelee.com" 10.1.239.157/healthz
curl --resolve cloudnative.jaquelee.com:443:10.1.239.157 https://cloudnative.jaquelee.com/healthz -k
