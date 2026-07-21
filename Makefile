# See: https://gist.github.com/asukakenji/f15ba7e588ac42795f421b48b8aede63
# For a list of valid GOOS and GOARCH values
# Note: these can be overriden on the command line e.g. `make PLATFORM=<platform> ARCH=<arch>`
PLATFORM?=$(shell go env GOOS)
ARCH?=$(shell go env GOARCH)

MAIN=cmd/helper-reset-password/main.go
DOCKER_IMAGE=portainer/helper-reset-password
ALL_OSVERSIONS.windows := 1809 1909 2004 20H2 ltsc2022
GOTESTSUM=go run gotest.tools/gotestsum@latest

ifeq ("$(PLATFORM)", "windows")
bin=helper-reset-password.exe
else
bin=helper-reset-password
endif

.PHONY: build pre clean lint test run
.PHONY: release-linux-amd64 release-linux-arm release-linux-arm64 release-windows-amd64
.PHONY: release-linux release-windows release manifest
.PHONY: build-linux-amd64 build-linux-arm build-linux-arm64 build-windows-amd64

# Standard build target that accepts PLATFORM and ARCH parameters
pre:
	@mkdir -p bin/

build: pre
	@echo "Building helper-reset-password for $(PLATFORM)/$(ARCH)..."
	GOOS=$(PLATFORM) GOARCH=$(ARCH) CGO_ENABLED=0 go build -o bin/$(bin) $(MAIN)
	@echo "✓ Build completed: bin/$(bin)"

run:
	go run $(MAIN)

# Legacy platform-specific targets (kept for backwards compatibility)
build-linux-amd64:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/helper-reset-password $(MAIN)

build-linux-arm:
	GOOS=linux GOARCH=arm CGO_ENABLED=0 go build -o bin/helper-reset-password $(MAIN)

build-linux-arm64:
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o bin/helper-reset-password $(MAIN)

build-linux-ppc64le:
	GOOS=linux GOARCH=ppc64le CGO_ENABLED=0 go build -o bin/helper-reset-password $(MAIN)

build-windows-amd64:
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o bin/helper-reset-password.exe $(MAIN)

release-linux-amd64: build-linux-amd64 image-linux-amd64
release-linux-arm: build-linux-arm image-linux-arm
release-linux-arm64: build-linux-arm64 image-linux-arm64
release-windows-amd64: build-windows-amd64 image-windows-amd64

release-linux: release-linux-amd64 release-linux-arm release-linux-arm64
release-windows: release-windows-amd64
release: release-linux release-windows manifest

image-linux-amd64:
	docker buildx build --output=type=registry --platform linux/amd64 -t $(DOCKER_IMAGE):linux-amd64 -f Dockerfile .

image-linux-arm:
	docker buildx build --output=type=registry --platform linux/arm/v7 -t $(DOCKER_IMAGE):linux-arm -f Dockerfile .

image-linux-arm64:
	docker buildx build --output=type=registry --platform linux/arm64 -t $(DOCKER_IMAGE):linux-arm64 -f Dockerfile .

# Use buildx to build Windows images
image-windows-amd64:
	for osversion in $(ALL_OSVERSIONS.windows); do \
		docker buildx build --output=type=registry --platform windows/amd64 -t $(DOCKER_IMAGE):windows$${osversion}-amd64 --build-arg OSVERSION=$${osversion} -f ./Dockerfile.windows . ; \
	done

clean:
	rm -rf bin/$(BINARY)*

manifest:
	manifest_image_folder=`echo "docker.io/$(DOCKER_IMAGE)" | sed "s|/|_|g" | sed "s/:/-/"`; \
	docker -D manifest create $(DOCKER_IMAGE):latest \
		$(DOCKER_IMAGE):linux-amd64 \
		$(DOCKER_IMAGE):linux-arm \
		$(DOCKER_IMAGE):linux-arm64 \
		$(DOCKER_IMAGE):windows1809-amd64 \
		$(DOCKER_IMAGE):windows1909-amd64 \
		$(DOCKER_IMAGE):windows2004-amd64 \
		$(DOCKER_IMAGE):windows20H2-amd64 \
		$(DOCKER_IMAGE):windowsltsc2022-amd64 ; \
	docker manifest annotate $(DOCKER_IMAGE):latest $(DOCKER_IMAGE):linux-arm --os linux --arch arm ; \
	docker manifest annotate $(DOCKER_IMAGE):latest $(DOCKER_IMAGE):linux-arm64 --os linux --arch arm64 ; \
	for osversion in $(ALL_OSVERSIONS.windows); do \
		BASEIMAGE=mcr.microsoft.com/windows/nanoserver:$${osversion} ; \
		full_version=`docker manifest inspect $${BASEIMAGE} | jq -r '.manifests[0].platform["os.version"]'`; \
		sed -i -r "s/(\"os\"\:\"windows\")/\0,\"os.version\":\"$${full_version}\"/" "$(DOCKER_CONFIG)/manifests/$${manifest_image_folder}-latest/$${manifest_image_folder}-windows$${osversion}-amd64" ; \
	done

	docker manifest push $(DOCKER_IMAGE):latest

lint:
	golangci-lint run --timeout=10m --new-from-rev=HEAD~ -c .golangci.yaml

test:
	$(GOTESTSUM) --format pkgname-and-test-fails --format-hide-empty-pkg --hide-summary skipped -- -cover -covermode=atomic -coverprofile=coverage.out ./...