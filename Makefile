.PHONY: test lint

test:
	go test -v ./...

lint:
	go -C tools tool golangci-lint run ../...
