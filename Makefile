MAIN=cmd/helper-reset-password/main.go
BINARY=helper-reset-password
DOCKER_IMAGE=portainer/helper-reset-password

build-linux: compile-linux-amd64 image-linux-amd64 compile-linux-arm image-linux-arm compile-linux-arm64 image-linux-arm64
build-windows: compile-windows-amd64 image-windows-amd64
release: build-linux build-windows manifest

run:
	go run $(MAIN)

compile-linux-amd64:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/helper-reset-password $(MAIN)

compile-linux-arm:
	GOOS=linux GOARCH=arm CGO_ENABLED=0 go build -o bin/helper-reset-password $(MAIN)

compile-linux-arm64:
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o bin/helper-reset-password $(MAIN)

compile-windows-amd64:
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o bin/helper-reset-password.exe $(MAIN)

image-linux-amd64:
	docker build -f Dockerfile -t $(DOCKER_IMAGE):linux-amd64 .; \
	docker push $(DOCKER_IMAGE):linux-amd64

image-linux-arm:
	docker build -f Dockerfile -t $(DOCKER_IMAGE):linux-arm .; \
	docker push $(DOCKER_IMAGE):linux-arm

image-linux-arm64:
	docker build -f Dockerfile -t $(DOCKER_IMAGE):linux-arm64 .; \
	docker push $(DOCKER_IMAGE):linux-arm64

image-windows-amd64:
	docker build -f Dockerfile.windows -t $(DOCKER_IMAGE):windows-amd64 .; \
	docker push $(DOCKER_IMAGE):windows-amd64

clean:
	rm -rf bin/$(BINARY)*

manifest:
	DOCKER_CLI_EXPERIMENTAL=enabled docker manifest create $(DOCKER_IMAGE):latest $(DOCKER_IMAGE):linux-amd64 $(DOCKER_IMAGE):linux-arm $(DOCKER_IMAGE):linux-arm64 $(DOCKER_IMAGE):windows-amd64; \
    DOCKER_CLI_EXPERIMENTAL=enabled docker manifest annotate $(DOCKER_IMAGE):latest $(DOCKER_IMAGE):linux-arm --os linux --arch arm; \
    DOCKER_CLI_EXPERIMENTAL=enabled docker manifest annotate $(DOCKER_IMAGE):latest $(DOCKER_IMAGE):linux-arm64 --os linux --arch arm64; \
    DOCKER_CLI_EXPERIMENTAL=enabled docker manifest push $(DOCKER_IMAGE):latest
