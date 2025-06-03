#!/bin/sh

if [ ! -z $dev   ]; then CompileDaemon -log-prefix=false -build="go build -o search-service ./openslides-search-service/cmd/searchd/main.go" -command="./search-service"; fi
if [ ! -z $tests ]; then go vet ./... && go test -test.short ./...; fi
if [ ! -z $prod  ]; then exec ./openslides-search-service; fi