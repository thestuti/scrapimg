FROM golang:1.18-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY src/ ./src/

COPY urls.txt /app/urls.txt

RUN go build -o extract-images ./src/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/extract-images /app/extract-images

RUN apk --no-cache add ca-certificates

ENTRYPOINT ["/app/extract-images"]
