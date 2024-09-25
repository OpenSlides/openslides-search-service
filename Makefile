build-dev:
	rm -fr openslides-autoupdate-service
	cp -r ../openslides-autoupdate-service .
	docker build . --target development --tag openslides-search-dev
	rm -fr openslides-autoupdate-service

run-tests:
	docker build . --target testing --tag openslides-search-test
	docker run openslides-search-test

all: gofmt gotest golinter

gotest:
	go test ./...

golinter:
	golint -set_exit_status ./...

gofmt:
	gofmt -l -s -w .
