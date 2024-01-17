
.EXPORT_ALL_VARIABLES:

BIN_DIR := ./bin
OUT_DIR := ./output
$(shell mkdir -p $(BIN_DIR) $(OUT_DIR))

APP_NAME?=imrenagicom-app
SERVER_NAME=imrenagicom-app
PACKAGE=github.com/imrenagicom/demo-app
TRACK?=stable
IMAGE_REGISTRY=imrenagi
IMAGE_NAME=$(IMAGE_REGISTRY)/app
IMAGE_TAG?=latest

CURRENT_DIR=$(shell pwd)
VERSION=$(shell cat ${CURRENT_DIR}/VERSION)
BUILD_DATE=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
GIT_COMMIT=$(shell git rev-parse HEAD)

DEPLOYMENT_TIMEOUT=600

STATIC_BUILD?=true

override LDFLAGS += \
  -X ${PACKAGE}.version=${VERSION} \
  -X ${PACKAGE}.buildDate=${BUILD_DATE}

ifeq (${STATIC_BUILD}, true)
override LDFLAGS += -extldflags "-static"
endif

ifneq (${GIT_TAG},)
LDFLAGS += -X ${PACKAGE}.gitTag=${GIT_TAG}
endif

BUF_VERSION:=v1.28.1
SWAGGER_UI_VERSION:=v4.15.5

bootstrap: install/protoc generate
	go install github.com/vektra/mockery/v2@v2.21.0
	go install github.com/golang/mock/mockgen@latest
	go install github.com/bufbuild/buf/cmd/buf@$(BUF_VERSION)
	go get ./...
	go mod tidy

generate: generate/proto generate/swagger-ui
	go generate ./...

install/protoc:
	go install \
            github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
            github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
            google.golang.org/protobuf/cmd/protoc-gen-go \
            google.golang.org/grpc/cmd/protoc-gen-go-grpc

generate/proto: install/protoc
	go run github.com/bufbuild/buf/cmd/buf@$(BUF_VERSION) generate

generate/swagger-ui:
	SWAGGER_UI_VERSION=$(SWAGGER_UI_VERSION) ./scripts/generate-swagger-ui.sh

binaries:
	CGO_ENABLED=0 GO111MODULE=on go build -a -ldflags '${LDFLAGS}' -o ${BIN_DIR}/course ./cmd/course/main.go

.PHONY: course/server
course/server:
	go run cmd/course/main.go server start --config course/conf/server.yaml --migration course/migrations

.PHONY: course/seed
course/seed:
	go run cmd/course/main.go server seed --config course/conf/server.yaml
