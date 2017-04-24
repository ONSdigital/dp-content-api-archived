
BUILD=build
HASH?=$(shell make hash)
BUILD_ARCH=$(BUILD)/$(GOOS)-$(GOARCH)
DATE:=$(shell date '+%Y%m%d-%H%M%S')
TGZ_FILE=dp-content-api-$(GOOS)-$(GOARCH)-$(DATE)-$(HASH).tar.gz

export GOOS?=$(shell go env GOOS)
export GOARCH?=$(shell go env GOARCH)

# The following variables are used to generate a nomad plan
export DATA_CENTER?=$(shell go env DATA_CENTER)
export S3_TAR_FILE?=$(shell go env S3_TAR_FILE)
export S3_CONTENT_URL?=$(shell go env S3_CONTENT_URL)
export S3_CONTENT_ACCESS_KEY?=$(shell go env S3_CONTENT_ACCESS_KEY)
export S3_CONTENT_SECRET_ACCESS_KEY?=$(shell go env S3_CONTENT_SECRET_ACCESS_KEY)
export S3_CONTENT_BUCKET?=$(shell go env S3_CONTENT_BUCKET)
export DATABASE_URL?=$(shell go env DATABASE_URL)
export DP_GENERATOR_URL?=$(shell go env DP_GENERATOR_URL)


build:
	@mkdir -p $(BUILD_ARCH)
	go build -o $(BUILD_ARCH)/bin/dp-content-api cmd/dp-content-api/main.go
	@cp -r static $(BUILD_ARCH)
	
package: build
	tar -zcf $(TGZ_FILE) -C $(BUILD_ARCH) .

nomad:
	@cp dp-content-api-template.nomad dp-content-api.nomad
	@sed -i.bak s,DATA_CENTER,$(DATA_CENTER),g dp-content-api.nomad
	@sed -i.bak s,S3_TAR_FILE,$(S3_TAR_FILE),g dp-content-api.nomad
	@sed -i.bak s,S3_CONTENT_URL,$(S3_CONTENT_URL),g dp-content-api.nomad
	@sed -i.bak s,S3_CONTENT_ACCESS_KEY,$(S3_CONTENT_ACCESS_KEY),g dp-content-api.nomad
	@sed -i.bak s,S3_CONTENT_SECRET_ACCESS_KEY,$(S3_CONTENT_SECRET_ACCESS_KEY),g dp-content-api.nomad
	@sed -i.bak s,S3_CONTENT_BUCKET,$(S3_CONTENT_BUCKET),g dp-content-api.nomad
	@sed -i.bak s,DATABASE_URL,$(DATABASE_URL),g dp-content-api.nomad
	@sed -i.bak s,DP_GENERATOR_URL,$(DP_GENERATOR_URL),g dp-content-api.nomad

hash:
	@git rev-parse --short HEAD

.PHONY: build package