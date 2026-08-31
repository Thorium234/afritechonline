.PHONY: test test-unit test-integration lint fmt

test: test-unit

test-unit:
	go test ./pkg/... ./internal/...

test-integration:
	go test -tags=integration ./...

lint:
	golangci-lint run

fmt:
	go fmt ./...
