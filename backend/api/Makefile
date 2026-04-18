.PHONY: proto/chat proto/social proto/notification

proto/chat:
	mkdir -p gen
	protoc \
	  --proto_path=proto \
	--go_out=gen \
	--go-grpc_out=gen \
	  proto/chat/chat.proto

proto/social:
	mkdir -p gen
	protoc \
	  --proto_path=proto \
	--go_out=gen \
	--go-grpc_out=gen \
	  proto/social_network/*.proto

proto/notification:
	mkdir -p gen
	protoc \
	  --proto_path=proto \
	--go_out=gen \
	--go-grpc_out=gen \
	  proto/notification/*.proto

proto-all: proto/chat proto/social proto/notification
