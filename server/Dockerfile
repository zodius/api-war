FROM golang:1.22-alpine AS builder
WORKDIR /app
RUN apk update && apk add build-base
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 go build -o server cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server /app/server
ENTRYPOINT ["/app/server"]