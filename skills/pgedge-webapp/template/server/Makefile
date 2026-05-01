.PHONY: build test coverage coverage-check lint vet fmt clean

build:
	cd src && go build -o ../../bin/<BINARY_NAME> ./cmd/<BINARY_NAME>

test:
	cd src && go test -race ./...

coverage:
	cd src && go test -race -covermode=atomic -coverprofile=coverage.out ./...
	cd src && go tool cover -func=coverage.out | tail -1

coverage-check: coverage
	cd src && go run ./tools/coverage-check \
		-profile coverage.out -threshold 90 \
		-exclude 'cmd/<BINARY_NAME>/main.go,cmd/<BINARY_NAME>/server.go,tools/coverage-check/main.go'

lint:
	cd src && golangci-lint run

vet:
	cd src && go vet ./...

fmt:
	cd src && gofmt -l -w .

clean:
	rm -f src/coverage.out
