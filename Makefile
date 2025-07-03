override SERVICE=search
override MAKEFILE_PATH=../dev/scripts/makefile
override DOCKER_COMPOSE_FILE=

# Build images for different contexts

build build-prod build-dev build-tests:
	bash $(MAKEFILE_PATH)/make-build-service.sh $@ $(SERVICE)

# Development

.PHONY: run-dev%

run-dev%:
	bash $(MAKEFILE_PATH)/make-run-dev.sh "$@" "$(SERVICE)" "$(DOCKER_COMPOSE_FILE)" "$(ARGS)" "$(USED_SHELL)"

# Tests

run-tests:
	bash dev/run-tests.sh

run-lint:
	gofmt -l -s -w .
	go test ./...
	golint -set_exit_status ./...


########################## Deprecation List ##########################

deprecation-warning:
	bash $(MAKEFILE_PATH)/make-deprecation-warning.sh

all:
	bash $(MAKEFILE_PATH)/make-deprecation-warning.sh "run-tests for tests and lints inside a container or run-lint for local linting"
	make gofmt
	make gotest
	make golinter

gotest: | deprecation-warning
	go test ./...

golinter: | deprecation-warning
	golint -set_exit_status ./...

gofmt: | deprecation-warning
	gofmt -l -s -w .
