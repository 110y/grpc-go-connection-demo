OS   := $(shell go env GOOS)
ARCH := $(shell go env GOARCH)

BUF_VERSION                := 1.0.0-rc5
PROTOC_GEN_GO_VERSION      := 1.27.1
PROTOC_GEN_GO_GRPC_VERSION := 1.1.0

BIN_DIR := $(shell pwd)/bin
TOOLS_DIR := ./tools
TOOLS_SUM := $(TOOLS_DIR)/go.sum

TOOLS_GO_BUILD_CMD := cd $(TOOLS_DIR) && go build -o

BUF                := $(abspath $(BIN_DIR)/buf)
PROTOC_GEN_GO      := $(abspath $(BIN_DIR)/protoc-gen-go)
PROTOC_GEN_GO_GRPC := $(abspath $(BIN_DIR)/protoc-gen-go-grpc)
CHANNELZCLI        := $(abspath $(BIN_DIR)/channelzcli)

buf: $(BUF)
$(BUF):
	curl -sSL "https://github.com/bufbuild/buf/releases/download/v${BUF_VERSION}/buf-$(shell uname -s)-$(shell uname -m)" -o $(BUF) && chmod +x $(BUF)

protoc-gen-go: $(PROTOC_GEN_GO)
$(PROTOC_GEN_GO):
	curl -sSL https://github.com/protocolbuffers/protobuf-go/releases/download/v$(PROTOC_GEN_GO_VERSION)/protoc-gen-go.v$(PROTOC_GEN_GO_VERSION).$(OS).$(ARCH).tar.gz | tar -C $(BIN_DIR) -xzv protoc-gen-go

protoc-gen-go-grpc: $(PROTOC_GEN_GO_GRPC)
$(PROTOC_GEN_GO_GRPC):
	curl -sSL https://github.com/grpc/grpc-go/releases/download/cmd%2Fprotoc-gen-go-grpc%2Fv$(PROTOC_GEN_GO_GRPC_VERSION)/protoc-gen-go-grpc.v$(PROTOC_GEN_GO_GRPC_VERSION).$(OS).$(ARCH).tar.gz | tar -C $(BIN_DIR) -xzv ./protoc-gen-go-grpc

channelzcli: $(CHANNELZCLI)
$(CHANNELZCLI): $(TOOLS_SUM)
	@$(TOOLS_GO_BUILD_CMD) $(CHANNELZCLI) github.com/kazegusuri/channelzcli

.PHONY: up
up:
	docker compose up --detach

.PHONY: pb
pb:
	docker run \
		--volume "$(shell pwd):/go/src/github.com/110y/grpc-go-connection-demo" \
		--workdir /go/src/github.com/110y/grpc-go-connection-demo \
		--rm \
		golang:1.17.1-bullseye \
		make gen-proto

.PHONY: gen-proto
gen-proto: $(BUF) $(PROTOC_GEN_GO) $(PROTOC_GEN_GO_GRPC)
	$(BUF) generate \
		--path ./caller/pb \
		--path ./callee/pb

.PHONY: show-caller-channel
show-caller-channel: $(CHANNELZCLI)
	 @$(CHANNELZCLI) -k --addr localhost:5000 ls channel
