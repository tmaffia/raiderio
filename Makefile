.PHONY: test lint setup

test:
	go test -v ./...

lint:
	go run -modfile=tools/go.mod github.com/golangci/golangci-lint/cmd/golangci-lint run ./...

setup:
	cd tools && go mod tidy