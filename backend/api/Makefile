.PHONY: proto/chat proto/social proto/notification proto-all clean

proto/chat:
	protoc \
	  --proto_path=proto \
	  --go_out=. \
	  --go_opt=module=api \
	  --go-grpc_out=. \
	  --go-grpc_opt=module=api \
	  proto/chat/chat.proto

proto/social:
	protoc \
	  --proto_path=proto \
	  --go_out=. \
	  --go_opt=module=api \
	  --go-grpc_out=. \
	  --go-grpc_opt=module=api \
	  proto/social_network/*.proto

proto/notification:
	protoc \
	  --proto_path=proto \
	  --go_out=. \
	  --go_opt=module=api \
	  --go-grpc_out=. \
	  --go-grpc_opt=module=api \
	  proto/notification/*.proto

proto-all: proto/chat proto/social proto/notification

clean:
	rm -rf gen
