.PHONY: protos

protots:
	#protoc -I testproto/ --go_out=. --go-grpc_out=. testproto/test.proto
	protoc -I testproto/ --go_out=. --go-grpc_out=require_unimplemented_servers=false:. testproto/test.proto
