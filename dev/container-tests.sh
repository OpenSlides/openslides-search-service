#!/bin/sh

CATCH=0

# Run Linters & Tests
go vet ./... || CATCH=1
go test -timeout 60s -race ./... || CATCH=1
gofmt -l . || CATCH=1
golint -set_exit_status ./... || CATCH=1

exit $CATCH
