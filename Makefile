
BUILD=build
HASH?=$(shell make hash)
BUILD_ARCH=$(BUILD)/$(GOOS)-$(GOARCH)
DATE:=$(shell date '+%Y%m%d-%H%M%S')
TGZ_FILE=dp-content-api-$(GOOS)-$(GOARCH)-$(DATE)-$(HASH).tar.gz

export GOOS?=$(shell go env GOOS)
export GOARCH?=$(shell go env GOARCH)

all: build nomad

build:
	@mkdir -p $(BUILD_ARCH)
	go build -o $(BUILD_ARCH)/bin/dp-content-api cmd/dp-content-api/main.go
	@cp -r static $(BUILD_ARCH)

test:
	go test -cover content/*.go
	
package: build
	tar -zcf $(TGZ_FILE) -C $(BUILD_ARCH) .

nomad:
	@for t in *-template.nomad; do			\
		plan=$${t%-template.nomad}.nomad;	\
		test -f $$plan && rm $$plan;		\
		sed	-e 's,DATA_CENTER,$(DATA_CENTER),g'			\
			-e 's,S3_TAR_FILE,$(S3_TAR_FILE),g'			\
			-e 's,S3_CONTENT_URL,$(S3_CONTENT_URL),g'		\
			-e 's,S3_CONTENT_ACCESS_KEY,$(S3_CONTENT_ACCESS_KEY),g'	\
			-e 's,S3_CONTENT_SECRET_ACCESS_KEY,$(S3_CONTENT_SECRET_ACCESS_KEY),g'	\
			-e 's,S3_CONTENT_BUCKET,$(S3_CONTENT_BUCKET),g'	\
			-e 's,DATABASE_URL,$(DATABASE_URL),g'		\
			-e 's,DP_GENERATOR_URL,$(DP_GENERATOR_URL),g'	\
			-e 's,HEALTHCHECK_ENDPOINT,$(HEALTHCHECK_ENDPOINT),g'	\
			< $$t > $$plan || exit 2;			\
	done

clean:
	test -d $(BUILD) && rm -r $(BUILD)
	for t in *-template.nomad; do plan=$${t%-template.nomad}.nomad; test -f $$plan || continue; rm $$plan || exit 2; done

hash:
	@git rev-parse --short HEAD

.PHONY: build package
