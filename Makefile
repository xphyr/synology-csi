# Copyright 2021 Synology Inc.

REGISTRY_NAME=xphyr
IMAGE_NAME=synology-csi
IMAGE_VERSION=v1.5.0
IMAGE_TAG=$(REGISTRY_NAME)/$(IMAGE_NAME):$(IMAGE_VERSION)

# For now, only build linux/amd64 platform
ifeq ($(GOARCH),)
GOARCH:=amd64
endif
GOARM?=""
BUILD_ENV=CGO_ENABLED=0 GOOS=linux GOARCH=$(GOARCH) GOARM=$(GOARM)
BUILD_FLAGS="-s -w -extldflags \"-static\""

.PHONY: all clean synology-csi-driver synocli test docker-build docker-build-multiarch vet lint fmt release goreleaser-docker

all: vet release

# Format all Go files
fmt:
	go fmt ./...

# Run go fmt, go vet, and golangci-lint (requires golangci-lint to be installed)
vet:
	$(MAKE) fmt
	golangci-lint run --timeout=5m

# Run only golangci-lint
lint:
	golangci-lint run --timeout=5m

release:
	goreleaser release --clean --skip-publish

# Build and publish docker images using GoReleaser
# Requires a properly configured .goreleaser.yaml

goreleaser-docker:
	goreleaser release --clean --skip-publish --skip-sign --rm-dist

docker-build:
	docker build -f Dockerfile -t $(IMAGE_TAG) .

docker-build-multiarch:
	docker buildx build -t $(IMAGE_TAG) --platform linux/amd64,linux/arm/v7,linux/arm64 . --push

synocli:
	@mkdir -p bin
	$(BUILD_ENV) go build -v -ldflags $(BUILD_FLAGS) -o ./bin/synocli ./synocli

test:
	go clean -testcache
	go test -v ./test/...

clean:
	-rm -rf ./bin

