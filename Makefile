run:
	go run ./webstreaming/cmd/api

generate:
	protoc --go_out=. --go-grpc_out=. \
    proto/service.proto
