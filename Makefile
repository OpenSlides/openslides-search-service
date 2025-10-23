override SERVICE=search

# Build images for different contexts

build-prod:
	docker build ./ $(ARGS) --tag "openslides-$(SERVICE)" --build-arg CONTEXT="prod" --target "prod"

build-dev:
	docker build ./ $(ARGS) --tag "openslides-$(SERVICE)-dev" --build-arg CONTEXT="dev" --target "dev"

build-tests:
	docker build ./ $(ARGS) --tag "openslides-$(SERVICE)-tests" --build-arg CONTEXT="tests" --target "tests"

# Tests

run-tests:
	bash dev/run-tests.sh

run-tests-local-branch:
	bash dev/run-tests.sh true

lint:
	bash dev/run-lint.sh -l

gofmt:
	gofmt -l -s -w .

# Local Development
run-clean-psql-setup:
	make -C .. dev-stop search compose-local-branch
	make -C .. dev-detached search compose-local-branch
	sleep 6
	make -C .. dev-exec search EXEC_COMMAND="search bash ./dev/create-models.sh true"
	sleep 6
	make curl-search-string-default

curl-search-string:
	curl "http://localhost:9050/system/search?q=$(Q)&c=$(C)"

curl-search-string-default:
	make curl-search-string Q=test C=topic
	make -C .. dev-log search compose-local-branch

log:
	make -C .. dev-log search compose-local-branch

########################## Deprecation List ##########################

deprecation-warning:
	@echo "\033[1;33m DEPRECATION WARNING: This make command is deprecated and will be removed soon! \033[0m"

deprecation-warning-alternative: | deprecation-warning
	@echo "\033[1;33m Please use the following command instead: $(ALTERNATIVE) \033[0m"

all:
	@make deprecation-warning-alternative ALTERNATIVE="run-tests for tests and lints inside a container or run-lint for local linting"
	make gofmt
	make gotest
	make golinter

gotest: | deprecation-warning
	go test ./...

golinter: | deprecation-warning
	golint -set_exit_status ./...
