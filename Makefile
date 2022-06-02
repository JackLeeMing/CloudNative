export tag=v1.0.0

root:
		export ROOT=github.com/JackLeeMing/CloudNative
build:
		echo "bulding package"
		mkdir -p bin/amd64
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/amd64 .
release:build
		echo "release package"
		docker build -t jackleeming/cloudnative:${tag} .
push:release
		echo "push package"
		docker push jackleeming/cloudnative:${tag}
