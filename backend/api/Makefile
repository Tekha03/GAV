.PHONY: proto/chat proto/social

proto/chat:
	protoc \
	  --proto_path=proto \
	  --go_out=chat_gen \
	  --go-grpc_out=chat_gen \
	  proto/chat/chat.proto

proto/social:
	protoc \
	  --proto_path=proto \
	  --go_out=social_gen \
	  --go-grpc_out=social_gen \
	  proto/social/social.proto

proto-all: proto/chat proto/social
