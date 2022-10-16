.PHONY: build
build:
	go build -o application -v ./cmd/main.go

.PHONY: run
run: build
	./application

.PHONY: compose-up
compose-up:
	docker-compose up

.PHONY: tidy
tidy:
	go mod tidy

.DEFAULT_GOAL := run
