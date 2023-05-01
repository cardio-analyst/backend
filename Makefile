.PHONY: compose-up
compose-up:
	docker-compose -f deployments/docker-compose.yml up

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: proto-auth
proto-auth:
	protoc -I api/proto/auth --go_out=api/proto/auth --go-grpc_out=api/proto/auth api/proto/auth/*.proto

.PHONY: proto-analytics
proto-analytics:
	protoc -I api/proto/analytics --go_out=api/proto/analytics --go-grpc_out=api/proto/analytics api/proto/analytics/*.proto

.PHONY: proto
proto: proto-auth proto-analytics

.DEFAULT_GOAL := compose-up
