.PHONY: build
build:
	go build -o application.exe -v ./cmd/main.go

.PHONY: run
run: build
	./application.exe

.PHONY: tidy
tidy:
	go mod tidy

.DEFAULT_GOAL := run
