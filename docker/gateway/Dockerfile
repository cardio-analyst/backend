FROM golang:1.17-buster as builder

WORKDIR /app/

COPY go.* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o application cmd/gateway/main.go

FROM alpine:3.15.4

WORKDIR /app/
COPY --from=builder /app/application ./
COPY --from=builder /app/configs/gateway/*.yaml ./
COPY --from=builder /app/assets/font ./assets/font
CMD ["/app/application", "-c", "/app/config.yaml"]
