BIN ?= main
DEST ?= ./bin
GO ?= go

.PHONY: all build build-arm build-arm64 clean fmt lint test bake image

all: build

build:
	$(GO) build -o $(BIN) ./main.go

build-arm:
	GOARCH=arm GOOS=linux $(GO) build -o $(BIN)-armv7 ./main.go

build-arm64:
	GOARCH=arm64 GOOS=linux $(GO) build -o $(BIN)-arm64 ./main.go

bake:
	docker buildx bake

image:
	docker buildx build --platform=linux/arm64,linux/arm/v7 --tag hello -o type=oci,dest=hello.tar .

fmt:
	$(GO) fmt ./...

lint:
	$(GO) vet ./...

test:
	$(GO) test -v ./...

clean:
	rm -f $(BIN) $(BIN)-armv7 $(BIN)-arm64 hello.tar
