#!/bin/bash

set -e

# Executes all tests. Should errors occur, CATCH will be set to 1, causing an erronoeus exit code.

echo "########################################################################"
echo "###################### Run Tests and Linters ###########################"
echo "########################################################################"

# Parameters
while getopts "s" FLAG; do
    case "${FLAG}" in
    s) SKIP_BUILD=true ;;
    *) echo "Can't parse flag ${FLAG}" && break ;;
    esac
done

# Setup
LOCAL_PWD=$(dirname "$0")

if [ -n "$1" ]
then
    COMPOSE_BRANCH=$(git -C "$SERVICE_FOLDER" branch --show-current)
else
    COMPOSE_BRANCH="main"
fi

# Helpers
USER_ID=$(id -u)
GROUP_ID=$(id -g)
DC="CONTEXT=tests USER_ID=$USER_ID GROUP_ID=$GROUP_ID COMPOSE_REFERENCE_BRANCH=$COMPOSE_BRANCH docker compose -f $LOCAL_PWD/../dev/docker-compose.dev.yml"

# Safe Exit
trap 'eval "$DC down --volumes"' EXIT INT TERM

# Execution
if [ -z "$SKIP_BUILD" ]; then make build-tests &> /dev/null; fi
eval "$DC up -d"

## Setup database
sleep 6
eval "$DC exec search bash dev/create-models.sh"

## Execute tests
eval "$DC exec search go test -timeout 60s -race ./..."


# Linters
bash "$LOCAL_PWD"/run-lint.sh -s
