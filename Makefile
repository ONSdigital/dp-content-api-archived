
BUILD=build
HASH?=$(shell make hash)
BUILD_ARCH=$(BUILD)/$(GOOS)-$(GOARCH)
DATE:=$(shell date '+%Y%m%d-%H%M%S')
TGZ_FILE=dp-content-api-$(GOOS)-$(GOARCH)-$(DATE)-$(HASH).tar.gz

export GOOS?=$(shell go env GOOS)
export GOARCH?=$(shell go env GOARCH)

build:
	@mkdir -p $(BUILD_ARCH)
	go build -o $(BUILD_ARCH)/bin/dp-content-api cmd/dp-content-api/main.go
	@cp -r static $(BUILD_ARCH)
	
package: build
	tar -zcf $(TGZ_FILE) -C $(BUILD_ARCH) .

hash:
	@git rev-parse --short HEAD

.PHONY: build package