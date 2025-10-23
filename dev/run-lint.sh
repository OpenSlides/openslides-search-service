#!/bin/bash

# Executes all linters. Should errors occur, CATCH will be set to 1, causing an erroneous exit code.

echo "########################################################################"
echo "###################### Run Linters #####################################"
echo "########################################################################"

# Parameters
while getopts "ls" FLAG; do
    case "${FLAG}" in
    l) LOCAL=true ;;
    s) SKIP_SETUP=true ;;
    *) echo "Can't parse flag ${FLAG}" && break ;;
    esac
done

# Setup
CONTAINER_NAME="search-tests"
IMAGE_TAG=openslides-search-tests
DOCKER_EXEC="docker exec ${CONTAINER_NAME}"

# Helpers
USER_ID=$(id -u)
GROUP_ID=$(id -g)
DC="CONTEXT=tests USER_ID=$USER_ID GROUP_ID=$GROUP_ID COMPOSE_REFERENCE_BRANCH=$COMPOSE_BRANCH docker compose -f $LOCAL_PWD/../dev/docker-compose.dev.yml"

# Safe Exit
trap 'if [ -z "$LOCAL" ] && [ -z "$SKIP_SETUP" ]; then eval "$DC down --volumes"; fi' EXIT

# Execution
if [ -z "$LOCAL" ]
then
    # Setup
    if [ -z "$SKIP_SETUP" ]
    then
        make build-tests >/dev/null 2>&1
        eval "$DC up -d"
    fi

    # Container Mode
    eval "$DC exec search go vet ./..."
    eval "$DC exec search golint -set_exit_status ./..."
    eval "$DC exec search gofmt -l ."
else
    # Local Mode
    go vet ./...
    golint -set_exit_status ./...
    gofmt -l -s -w .
fi
