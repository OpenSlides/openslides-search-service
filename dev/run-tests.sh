#!/bin/bash

# Executes all tests. Should errors occur, CATCH will be set to 1, causing an erronoeus exit code.

echo "########################################################################"
echo "###################### Run Tests and Linters ###########################"
echo "########################################################################"

# Setup
IMAGE_TAG=openslides-search-tests
LOCAL_PWD=$(dirname "$0")

if [ -n "$1" ]
then
    COMPOSE_BRANCH=$(git -C "$SERVICE_FOLDER" branch --show-current)
else
    COMPOSE_BRANCH="main"
fi

# Safe Exit
#trap 'eval "CONTEXT=tests USER_ID=1000 GROUP_ID=1000 COMPOSE_REFERENCE_BRANCH=$COMPOSE_BRANCH docker compose -f $LOCAL_PWD/../dev/docker-compose.dev.yml down --volumes"' EXIT INT TERM

# Execution

make build-test
eval "CONTEXT=tests USER_ID=1000 GROUP_ID=1000 COMPOSE_REFERENCE_BRANCH=$COMPOSE_BRANCH docker compose -f $LOCAL_PWD/../dev/docker-compose.dev.yml up -d"

## Setup database
sleep 6

eval "CONTEXT=tests USER_ID=1000 GROUP_ID=1000 COMPOSE_REFERENCE_BRANCH=$COMPOSE_BRANCH docker compose -f $LOCAL_PWD/../dev/docker-compose.dev.yml exec search bash create-models.sh"

sleep 6

## Execute tests
eval "CONTEXT=tests USER_ID=1000 GROUP_ID=1000 COMPOSE_REFERENCE_BRANCH=$COMPOSE_BRANCH docker compose -f $LOCAL_PWD/../dev/docker-compose.dev.yml exec search go vet ./... && go test -timeout 60s -race ./... && gofmt -l . && golint -set_exit_status ./..."
