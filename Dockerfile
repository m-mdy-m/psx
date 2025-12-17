FROM golang:1.25.0-alpine AS builder
RUN apk add --no-cache git make
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w -X main.Version=$(git describe --tags --always --dirty 2>/dev/null || echo 'docker')" \
    -o psx \
    ./cmd/psx
RUN chmod +x psx && ./psx --version
FROM alpine:3.19
RUN apk add --no-cache ca-certificates
RUN addgroup -g 1000 psx && \
    adduser -D -u 1000 -G psx psx
WORKDIR /project
COPY --from=builder /build/psx /usr/local/bin/psx
RUN chown psx:psx /usr/local/bin/psx
USER psx
ENTRYPOINT ["psx"]
CMD ["--help"]
LABEL org.opencontainers.image.title="PSX - Project Structure Checker"
LABEL org.opencontainers.image.description="Validate and standardize project structures"
LABEL org.opencontainers.image.authors="Genix <bitsgenix@gmail.com>"
LABEL org.opencontainers.image.source="https://github.com/m-mdy-m/psx"
LABEL org.opencontainers.image.documentation="https://github.com/m-mdy-m/psx"
LABEL org.opencontainers.image.licenses="MIT"