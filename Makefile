SHELL=/bin/bash

# Set this -->[/Users/xxxx/Sonr/]<-- to Folder of Sonr Repos
SONR_ROOT_DIR=/Users/prad/Developer
ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

# Set this -->[/Users/xxxx/Sonr/]<-- to Folder of Sonr Repos
PROTO_DEF_PATH=/Users/prad/Developer/core/proto

# @ Proto Directories
PROTO_LIST_ALL=${ROOT_DIR}/proto/**/*.proto
MODULE_NAME=github.com/sonr-io/core
GO_OPT_FLAG=--go_opt=module=${MODULE_NAME}
GRPC_OPT_FLAG=--go-grpc_opt=module=${MODULE_NAME}
PROTO_GEN_GO="--go_out=."
PROTO_GEN_RPC="--go-grpc_out=."
PROTO_GEN_DOCS="--doc_out=docs"

all: Makefile
	@figlet -f larry3d Sonr Core
	@echo ''
	@sed -n 's/^##//p ' $<

docker:
	@echo "----"
	@echo "Sonr: Building and Pushing Docker Image"
	@echo "----"
	@docker build . -t ghcr.io/sonr-io/highway
	@docker push ghcr.io/sonr-io/highway:latest

## [protobuf]     :   Compiles Protobuf models for Core Library and Plugin
protobuf:
	@echo "----"
	@echo "Sonr: Compiling Protobufs"
	@echo "----"
	@echo "Generating Protobuf Go code..."
	@protoc $(PROTO_LIST_ALL) --proto_path=$(ROOT_DIR) $(PROTO_GEN_GO) $(GO_OPT_FLAG)
	@protoc $(PROTO_LIST_ALL) --proto_path=$(ROOT_DIR) $(PROTO_GEN_RPC) $(GRPC_OPT_FLAG)

## [clean]     :   Reinitializes Gomobile and Removes Framworks from Plugin
clean:
	cd $(CORE_BIND_DIR) && $(GOCLEAN)
	go mod tidy
	go clean -cache -x
	rm -rf $(BIND_DIR_IOS)
	rm -rf $(BIND_DIR_ANDROID)
	mkdir -p $(BIND_DIR_IOS)
	mkdir -p $(BIND_DIR_ANDROID)
	cd $(CORE_BIND_DIR) && gomobile init
