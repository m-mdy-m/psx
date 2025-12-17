# builder stage
FROM golang:1.25.0-alpine AS builder
RUN apk add --no-cache git make
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w -X main.Version=$(git describe --tags --always --dirty 2>/dev/null || echo 'docker')" \
    -o /usr/local/bin/psx \
    ./cmd/psx

# final image
FROM alpine:3.18
RUN apk add --no-cache ca-certificates
COPY --from=builder /usr/local/bin/psx /usr/local/bin/psx
RUN chmod +x /usr/local/bin/psx
WORKDIR /project
ENTRYPOINT ["/usr/local/bin/psx"]
