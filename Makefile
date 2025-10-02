SERVICE=search

build-prod:
	docker build ./ --tag "openslides-$(SERVICE)" --build-arg CONTEXT="prod" --target "prod"

build-dev:
	docker build ./ --tag "openslides-$(SERVICE)-dev" --build-arg CONTEXT="dev" --target "dev"

build-test:
	docker build ./ --tag "openslides-$(SERVICE)-tests" --build-arg CONTEXT="tests" --target "tests"

run-tests:
	bash dev/run-tests.sh

run-tests-local-branch:
	bash dev/run-tests.sh true

all: gofmt gotest golinter

gotest:
	go test ./...

golinter:
	golint -set_exit_status ./...

gofmt:
	gofmt -l -s -w .

# Local Development
run-clean-psql-setup:
	make -C .. dev-stop search compose-local-branch
	make -C .. dev-detached search compose-local-branch
	sleep 6
	make -C .. dev-exec search EXEC_COMMAND="search bash create-models.sh true"
	sleep 6
	make curl-search-string-default

curl-search-string:
	curl "http://localhost:9050/system/search?q=$(Q)&c=$(C)"

curl-search-string-default:
	make curl-search-string Q=test C=topic
	make -C .. dev-log search compose-local-branch

log:
	make -C .. dev-log search compose-local-branch
