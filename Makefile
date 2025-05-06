SHELL=/bin/bash
# make 指令
MAKE=make
# Go parameters
GO_CMD=go
GO_BUILD=$(GO_CMD) build
GO_CLEAN=$(GO_CMD) clean
GO_TEST=$(GO_CMD) test
GO_GET=$(GO_CMD) get
GO_VET=$(GO_CMD) vet
GO_RUN=$(GO_CMD) run
#測試
COVER_PROFILE=cover.out

#docker
DOCKER_CMD=docker
DOCKER_BUILD=$(DOCKER_CMD) build
DOCKER_PUSH=$(DOCKER_CMD) push
DOCKER_IMAGE_REGISTRY=registry.digiwincloud.com.cn/ops
DOCKER_IMAGE_NAME=ops
DOCKER_FULL_IMAGE=$(DOCKER_IMAGE_REGISTRY)$(DOCKER_IMAGE_NAME):$(VERSION).$(shell cat $(SUB_VERSION_FILE))
#打包
BINARY_NAME=audit
BINARY_UNIX=$(BINARY_NAME)_unix
VERSION:=$(shell cat VERSION)
#版本控制
SUB_VERSION_FILE=./version_control/BUILD

#Git
GIT_CMD=git
GIT_BRANCH=$(GIT_CMD) branch
GIT_ADD=$(GIT_CMD) add
GIT_COMMIT=$(GIT_CMD) commit
GIT_PUSH=$(GIT_CMD) push
GIT_CURRENT_BRANCH=$(GIT_BRANCH) --show-current


#其他指令
ALL_PATH=./...

all: deps test build pack
deps:
	$(GO_GET) -u $(ALL_PATH)
test:
	$(GO_TEST) -v $(ALL_PATH) -cover
test_output:
	$(GO_TEST) -v $(ALL_PATH) -coverprofile=$(COVER_PROFILE)
build:
	$(GO_BUILD) -o $(BINARY_NAME)
run:
	$(GO_RUN) main.go
clean:
	$(GO_CLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
pack:
	tar -cvzf $(BINARY_NAME)-v$(VERSION).tar.gz $(BINARY_NAME)
docker_build:
	@echo "開始打包 Docker Image - $(DOCKER_FULL_IMAGE)"
	$(DOCKER_BUILD) -t $(DOCKER_FULL_IMAGE) .
docker_push:
	@echo "開始 push docker image - $(DOCKER_FULL_IMAGE)"
	$(DOCKER_PUSH) $(DOCKER_FULL_IMAGE)
docker_ci: vc docker_build docker_push add_tag

vc:
	@make -C version_control
add_tag:
	@make -C version_control commit_record add_tag