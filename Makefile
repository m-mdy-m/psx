VERSION= ${shell git describe --tags --always --dirty}

LDFLAGS= -ldflags "-X main.Version=$(VERSION)"

build:
	@echo "Building PSX $(VERSION)..."
	go build $(LDFLAGS) -o build/psx ./cmd/psx
dev:
	@echo "Running in dev mode"
	go run ./cmd/psx

clean:
	rn -rf build/

