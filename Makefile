.PHONY: build clean deploy
GOBUILD=env GOOS=linux go build -ldflags="-s -w" -o
functions := $(shell find functions -name \*main.go | awk -F'/' '{print $$2}')

deploy: ensure-deps fmt linter clean build test
	@echo "==> Running deploy..."
	sls deploy --verbose

build:
	@echo "==> Running build..."
	rm -f bin/*
	@for function in $(functions) ; do \
		env GOOS=linux go build -ldflags="-s -w" -o bin/$$function functions/$$function/main.go ; \
	done

clean:
	@echo "==> Running clean..."
	rm -rf ./bin

test:
	@echo "==> Running tests..."
	go test ./...

### Formatting, linting, and deps
fmt:
	@echo "==> Running format..."
	go fmt ./...

linter:
	@echo "==> Running linter..."
	golangci-lint run ./...

ensure-deps:
	@echo "=> Syncing dependencies with go mod tidy"
	@go mod tidy


