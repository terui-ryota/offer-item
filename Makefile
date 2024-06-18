# backend-migration.mk
# バックエンドのマイグレ用の共通Makefile
GO ?= go

SQL_MIGRATE_VER := 1.2.0
MIGRATE_DIR ?= ./deployment/migration

.PHONY: deps/sql-migrate
deps/sql-migrate:
	@$(GO) install github.com/rubenv/sql-migrate/sql-migrate@v$(SQL_MIGRATE_VER)

.PHONY: migrate-status
migrate-status: deps/sql-migrate
migrate-status: ## sql-migrate statusを実行
	$(call exec-sql-migrate, status)

.PHONY: migrate-down
migrate-down: deps/sql-migrate
migrate-down: ## sql-migrate downを実行
	$(call exec-sql-migrate, down)

.PHONY: migrate-up
migrate-up: deps/sql-migrate
migrate-up: ## sql-migrate upを実行
	$(call exec-sql-migrate, up)

.PHONY: migrate-new
migrate-new: deps/sql-migrate
migrate-new: ## sql-migrate newを実行
	$(call exec-sql-migrate, new)

define exec-sql-migrate
	cd $(MIGRATE_DIR) && sql-migrate $(1) -config=dbconfig.yaml -env=local
endef


# -------------
# backend-common.mk
CGO_ENABLED ?= 0
BINDIR ?= bin
RELEASE ?= 0

#SERVICE_NAME := pick
COMPONENT_NAME ?=
#LIBGO_CONST_PACKAGE := github.com/ca-media-nantes/libgo/v2/consts
#GOLIB_CONST_PACKAGE := github.com/ca-media-nantes/pick/go-lib/pkg/consts
ROOT_DIR ?= $(shell git rev-parse --show-toplevel)

BINARIES:=  $(patsubst ./cmd/%/, $(BINDIR)/%, $(dir $(wildcard ./cmd/*/main.go)))

# build flags
GO_BUILD_TAGS ?=
GO_BUILD_FLAGS := -installsuffix netgo -tags=netgo,$(GO_BUILD_TAGS)
ifeq ($(RELEASE), 1)
	GO_BUILD_FLAGS += -a -trimpath
endif

ifdef ADDITIONAL_GO_BUILD_FLAGS
	GO_BUILD_FLAGS += $(ADDITIONAL_GO_BUILD_FLAGS)
endif

# 現状失敗する
.PHONY: build
build: ## ビルドを実行
build: $(BINARIES)

# 実ビルドタスク
$(BINARIES): GOOS ?= $(shell go env GOOS)
$(BINARIES): GOARCH ?= $(shell go env GOARCH)
$(BINARIES): VERSION ?= 0.0.0
$(BINARIES): VCS_REVISION ?= $(shell git rev-parse --short HEAD)
#$(BINARIES): GO_LDFLAGS += -X $(LIBGO_CONST_PACKAGE).ModuleName=$(COMPONENT_NAME)-$(patsubst $(COMPONENT_NAME)-%,%,$(patsubst $(BINDIR)/%,%,$@))
#$(BINARIES): GO_LDFLAGS += -X $(LIBGO_CONST_PACKAGE).VcsRevision=$(VCS_REVISION)
#$(BINARIES): GO_LDFLAGS += -X $(LIBGO_CONST_PACKAGE).Version=$(VERSION)
#$(BINARIES): FORCE
#	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO) build -o $@ $(GO_BUILD_FLAGS) -ldflags "$(GO_LDFLAGS)" $(@:$(BINDIR)/%=./cmd/%/)
## 強制的にビルドを実行するためのダミーターゲット
#.PHONY: FORCE
#FORCE:

SQLBOILER_CONFIG ?= build/sqlboiler/sqlboiler.toml
ENTITY_DIST_DIR ?= internal/pkg/infrastructure/db/entity

.PHONY: generate-entity
generate-entity: ## sqlboilerを実行
generate-entity: $(BINDIR)/sqlboiler $(BINDIR)/sqlboiler-mysql
	$(BINDIR)/sqlboiler $(BINDIR)/sqlboiler-mysql -c $(SQLBOILER_CONFIG) -o $(ENTITY_DIST_DIR) -p entity --no-tests --wipe

$(BINDIR)/sqlboiler:  go.mod go.sum
	GOBIN=$(abspath $(BINDIR)) $(GO) install github.com/volatiletech/sqlboiler

$(BINDIR)/sqlboiler-mysql: go.mod go.sum
	GOBIN=$(abspath $(BINDIR)) $(GO) install github.com/volatiletech/sqlboiler/drivers/sqlboiler-mysql


AWS_ECR_HOST ?= 192338233093.dkr.ecr.ap-northeast-1.amazonaws.com
.PHONY: build-image
ifdef PRE_BUILD_IMAGE_TARGETS
build-image: $(PRE_BUILD_IMAGE_TARGETS)
endif
build-image: TAG_NAME ?= latest
build-image: DOCKER_BUILD_ARGS ?=
build-image: ECR_IMAGE_NAME := $(AWS_ECR_HOST)/pick-$(COMPONENT_NAME):$(TAG_NAME)
build-image:  ## build image
	docker build -t backend-base -f $(ROOT_DIR)/docker/backend-base/Dockerfile $(ROOT_DIR)/docker/backend-base
	docker build -t $(ECR_IMAGE_NAME) $(patsubst %, --build-arg=%, $(DOCKER_BUILD_ARGS)) .
ifeq ($(CI), true)
	echo "build-images=$(ECR_IMAGE_NAME)" >> $$GITHUB_OUTPUT
endif

# -------------
# golang.mk

GO_TEST_OPTS := -v

.PHONY: generate
generate:  ## go generateを実行
	$(GO) generate -v -tags=wireinject ./...

.PHONY: test
test: ## テストを実行
	$(GO) test $(GO_TEST_OPTS) ./... $(PIPE_GO_TEST_RESULT)
