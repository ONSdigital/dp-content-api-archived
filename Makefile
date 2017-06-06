MAIN=dp-content-api
SHELL=bash

BUILD=build
HASH?=$(shell make hash)
BUILD_ARCH=$(BUILD)/$(GOOS)-$(GOARCH)
BIN_DIR?=.

DATE:=$(shell date '+%Y%m%d-%H%M%S')
TGZ_FILE=dp-content-api-$(GOOS)-$(GOARCH)-$(DATE)-$(HASH).tar.gz

PORT?=8082
HEALTHCHECK_ENDPOINT?=/healthcheck
DEV?=
S3_URL?=s3.amazonaws.com
S3_REGION?=eu-west-1

NOMAD?=
NOMAD_SRC_DIR?=nomad
NOMAD_PLAN_TARGET?=$(BUILD)
NOMAD_PLAN=$(NOMAD_PLAN_TARGET)/$(MAIN).nomad

export GOOS?=$(shell go env GOOS)
export GOARCH?=$(shell go env GOARCH)
thisOS:=$(shell uname -s)

ifeq ($(thisOS),Darwin)
SED?=gsed
else
SED?=sed
endif

ifdef DEV
DATA_CENTER?=dc1
HUMAN_LOG?=1
else
DATA_CENTER?=$(S3_REGION)
HUMAN_LOG?=
endif

all: build nomad

build:
	@mkdir -p $(BUILD_ARCH)
	go build -o $(BUILD_ARCH)/$(BIN_DIR)/dp-content-api cmd/dp-content-api/main.go
	@cp -r static $(BUILD_ARCH)

$(MAIN) run:
ifdef NOMAD
	@if [[ ! -f $(NOMAD_PLAN) ]]; then echo Cannot see $(NOMAD_PLAN); exit 1; fi; echo nomad run $(NOMAD_PLAN); nomad run $(NOMAD_PLAN)
else
	@main=$(CMD_DIR)/$@/main.go; if [[ ! -f $$main ]]; then echo Cannot see $$main; exit 1; fi; go run -race $$main
endif

test:
	go test -cover content/*.go
	
package: build
	tar -zcf $(TGZ_FILE) -C $(BUILD_ARCH) .

nomad:
	@test -d $(NOMAD_PLAN_TARGET) || mkdir -p $(NOMAD_PLAN_TARGET)
	@driver=exec; [[ -n "$(DEV)" ]] && driver=raw_exec;	\
	for t in *-template.nomad; do			\
		plan=$(NOMAD_PLAN_TARGET)/$${t%-template.nomad}.nomad;	\
		test -f $$plan && rm $$plan;		\
		$(SED) -r	\
			-e 's,\bDATA_CENTER\b,$(DATA_CENTER),g'			\
			-e 's,\bS3_TAR_FILE\b,$(S3_TAR_FILE),g'			\
			-e 's,\bS3_CONTENT_URL\b,$(S3_URL),g'			\
			-e 's,\bS3_CONTENT_BUCKET\b,$(S3_CONTENT_BUCKET),g'	\
			-e 's,\bDATABASE_URL\b,$(DATABASE_URL),g'		\
			-e 's,\bDP_GENERATOR_URL\b,$(DP_GENERATOR_URL),g'	\
			-e 's,\bHEALTHCHECK_ENDPOINT\b,$(HEALTHCHECK_ENDPOINT),g'	\
			-e 's,\bHUMAN_LOG_FLAG\b,$(HUMAN_LOG),g'			\
			-e 's,\bCONTENT_API_PORT\b,$(PORT),g'			\
			-e 's,^(  *driver  *=  *)"exec",\1"'$$driver'",'	\
			< $$t > $$plan || exit 2;			\
	done

clean:
	test -d $(BUILD) && rm -r $(BUILD)

hash:
	@git rev-parse --short HEAD

.PHONY: build package
