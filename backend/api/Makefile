.PHONY: proto/chat proto/social proto/notification proto-all clean

ROOT := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
PROTO_DIR := $(ROOT)proto
GEN_DIR := $(ROOT)gen

proto/chat:
	mkdir -p $(GEN_DIR)
	protoc \
	  --proto_path=$(PROTO_DIR) \
	  --go_out=$(GEN_DIR) \
	  --go-grpc_out=$(GEN_DIR) \
	  $(PROTO_DIR)/chat/chat.proto

proto/social:
	mkdir -p $(GEN_DIR)
	protoc \
	  --proto_path=$(PROTO_DIR) \
	  --go_out=$(GEN_DIR) \
	  --go-grpc_out=$(GEN_DIR) \
	  $(PROTO_DIR)/social_network/*.proto

proto/notification:
	mkdir -p $(GEN_DIR)
	protoc \
	  --proto_path=$(PROTO_DIR) \
	  --go_out=$(GEN_DIR) \
	  --go-grpc_out=$(GEN_DIR) \
	  $(PROTO_DIR)/notification/*.proto

proto-all: proto/chat proto/social proto/notification

clean:
	rm -rf $(GEN_DIR)
