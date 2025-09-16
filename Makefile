SERVICE=search

build-prod:
	docker build ./ --tag "openslides-$(SERVICE)" --build-arg CONTEXT="prod" --target "prod"

build-dev:
	docker build ./ --tag "openslides-$(SERVICE)-dev" --build-arg CONTEXT="dev" --target "dev"

build-test:
	docker build ./ --tag "openslides-$(SERVICE)-tests" --build-arg CONTEXT="tests" --target "tests"

run-tests:
	bash dev/run-tests.sh

all: gofmt gotest golinter

gotest:
	go test ./...

golinter:
	golint -set_exit_status ./...

gofmt:
	gofmt -l -s -w .

run-clean-psql-setup:
	make -C .. dev-stop search compose-local-branch
	make -C .. dev-detached search compose-local-branch
	make -C .. dev-enter search compose-local-branch ATTACH_TARGET_CONTAINER=search
