.PHONY: compose-up
compose-up:
	docker-compose up

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: proto-auth
proto-auth:
	protoc -I pkg/api/proto/auth --go_out=pkg/api/proto/auth --go-grpc_out=pkg/api/proto/auth pkg/api/proto/auth/*.proto

.PHONY: proto
proto: proto-auth

.DEFAULT_GOAL := compose-up
