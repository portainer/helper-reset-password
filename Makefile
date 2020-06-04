MAIN=cmd/helper-reset-password/main.go
BINARY=helper-reset-password
DOCKER_IMAGE=portainer/helper-reset-password

release-linux-amd64: build-linux-amd64 image-linux-amd64
release-linux-arm: build-linux-arm image-linux-arm
release-linux-arm64: build-linux-arm64 image-linux-arm64
release-windows-amd64: build-windows-amd64 image-windows-amd64

release-linux: release-linux-amd64 release-linux-arm release-linux-arm64
release-windows: release-windows-amd64
release: release-linux release-windows manifest

run:
	go run $(MAIN)

build-linux-amd64:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/helper-reset-password $(MAIN)

build-linux-arm:
	GOOS=linux GOARCH=arm CGO_ENABLED=0 go build -o bin/helper-reset-password $(MAIN)

build-linux-arm64:
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o bin/helper-reset-password $(MAIN)

build-windows-amd64:
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

# Requires a properly configured rebase-docker-image tool
image-windows-amd64:
	docker build -f Dockerfile.windows -t $(DOCKER_IMAGE):windows-amd64 .; \
	docker push $(DOCKER_IMAGE):windows-amd64; \
	rebase-docker-image $(DOCKER_IMAGE):windows-amd64 -t $(DOCKER_IMAGE):windows1809-amd64 -b "mcr.microsoft.com/windows/nanoserver:1809"; \
    rebase-docker-image $(DOCKER_IMAGE):windows-amd64 -t $(DOCKER_IMAGE):windows1903-amd64 -b "mcr.microsoft.com/windows/nanoserver:1903"

clean:
	rm -rf bin/$(BINARY)*

manifest:
	DOCKER_CLI_EXPERIMENTAL=enabled docker manifest create $(DOCKER_IMAGE):latest $(DOCKER_IMAGE):linux-amd64 $(DOCKER_IMAGE):linux-arm $(DOCKER_IMAGE):linux-arm64 $(DOCKER_IMAGE):windows1809-amd64 $(DOCKER_IMAGE):windows1903-amd64; \
    DOCKER_CLI_EXPERIMENTAL=enabled docker manifest annotate $(DOCKER_IMAGE):latest $(DOCKER_IMAGE):linux-arm --os linux --arch arm; \
    DOCKER_CLI_EXPERIMENTAL=enabled docker manifest annotate $(DOCKER_IMAGE):latest $(DOCKER_IMAGE):linux-arm64 --os linux --arch arm64; \
    DOCKER_CLI_EXPERIMENTAL=enabled docker manifest push -p $(DOCKER_IMAGE):latest
