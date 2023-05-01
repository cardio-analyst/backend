FROM golang:1.19-buster as builder

WORKDIR /app/

COPY go.* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o application cmd/analytics/main.go

FROM alpine:3.15.4

WORKDIR /app/
COPY --from=builder /app/application ./
COPY --from=builder /app/configs/analytics/*.yaml ./
CMD ["/app/application", "-c", "/app/config.yaml"]
