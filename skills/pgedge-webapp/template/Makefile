.PHONY: build build-server build-client test test-server test-client coverage \
        coverage-check lint format dev docker helm-lint helm-template clean openapi

build: build-server build-client

build-server:
	$(MAKE) -C server build

build-client:
	cd client && npm run build

test: test-server test-client

test-server:
	$(MAKE) -C server test

test-client:
	cd client && npm test

coverage:
	$(MAKE) -C server coverage
	cd client && npm run test:coverage

coverage-check:
	$(MAKE) -C server coverage-check
	cd client && npm run test:coverage

lint:
	$(MAKE) -C server lint
	cd client && npm run lint

format:
	$(MAKE) -C server fmt
	cd client && npm run format

dev:
	bash start_dev.sh

docker:
	docker build -t <BINARY_NAME>:dev -f docker/Dockerfile .

helm-lint:
	helm lint deploy/helm/<BINARY_NAME>

helm-template:
	helm template <BINARY_NAME> deploy/helm/<BINARY_NAME>

clean:
	rm -rf bin client/dist client/coverage server/src/coverage.out

# Regenerate docs/docs/api/openapi.json from the server's programmatic
# spec builder. Delegates to the server Makefile, which knows the
# binary path and runs `<binary> --generate-openapi-spec`.
openapi:
	$(MAKE) -C server openapi
