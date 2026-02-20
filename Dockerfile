ARG CONTEXT=prod

FROM golang:1.26.0-alpine AS base

## Setup
ARG CONTEXT
WORKDIR /app/openslides-search-service
ENV APP_CONTEXT=${CONTEXT}

## Installs
RUN apk add git --no-cache

COPY go.mod go.sum ./
RUN go mod download

COPY cmd cmd
COPY pkg pkg

## External Information
EXPOSE 9050

# Development Image
FROM base AS dev

## Installs
RUN ["go", "install", "github.com/githubnemo/CompileDaemon@latest"]
RUN apk add make bash-completion postgresql-client

COPY entrypoint.sh ./
COPY meta/search.yml ./
COPY meta/models.yml ./
COPY meta ./meta
COPY dev/mock_data.sql ./dev/mock_data.sql
COPY dev/create-models.sh ./dev/create-models.sh


## Entrypoint
ENTRYPOINT ["./entrypoint.sh"]

HEALTHCHECK CMD wget --spider -q http://localhost:9050/system/search/health || exit 1

## Command
CMD CompileDaemon -log-prefix=false -build="go build -o openslides-search-service ./cmd/searchd/main.go" -command="./openslides-search-service"

# Testing Image
FROM dev AS tests

COPY dev/container-tests.sh ./dev/container-tests.sh

RUN apk add --no-cache \
    build-base \
    docker && \
    go get -u github.com/ory/dockertest/v3 && \
    go install golang.org/x/lint/golint@latest && \
    chmod +x dev/container-tests.sh

## Command
STOPSIGNAL SIGKILL
CMD ["sleep", "inf"]

# Production Image
FROM base AS builder
RUN go build -o openslides-search-service ./cmd/searchd/main.go

FROM alpine:3 AS prod

## Setup
ARG CONTEXT
ENV APP_CONTEXT=prod

COPY entrypoint.sh /
COPY meta/search.yml /
COPY meta/models.yml /
COPY --from=builder /app/openslides-search-service/openslides-search-service /

## External Information
LABEL org.opencontainers.image.title="OpenSlides Search Service"
LABEL org.opencontainers.image.description="The Search Service is a http endpoint where the clients can search for data within Openslides."
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.source="https://github.com/OpenSlides/openslides-search-service"

EXPOSE 9050

## Command
ENTRYPOINT ["./entrypoint.sh"]

CMD exec ./openslides-search-service

HEALTHCHECK CMD wget --spider -q http://localhost:9050/system/search/health || exit 1
