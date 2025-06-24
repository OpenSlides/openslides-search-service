ARG CONTEXT=prod

FROM golang:1.24.3-alpine as base

## Setup
ARG CONTEXT
WORKDIR /app/openslides-search-service
ENV ${CONTEXT}=1


## Installs
RUN apk add git --no-cache

COPY go.mod go.sum ./
RUN go mod download

COPY cmd cmd
COPY pkg pkg

## External Information
LABEL org.opencontainers.image.title="OpenSlides Search Service"
LABEL org.opencontainers.image.description="The Search Service is a http endpoint where the clients can search for data within Openslides."
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.source="https://github.com/OpenSlides/openslides-search-service"

EXPOSE 9050




# Development Image
FROM base as dev

RUN ["go", "install", "github.com/githubnemo/CompileDaemon@latest"]

COPY entrypoint.sh ./
COPY meta/search.yml .
COPY meta/models.yml .

## Entrypoint
ENTRYPOINT ["./entrypoint.sh"]

## Command
CMD CompileDaemon -log-prefix=false -build="go build -o openslides-search-service ./cmd/searchd/main.go" -command="./openslides-search-service"


# Testing Image
FROM base as tests

RUN apk add build-base --no-cache

CMD go vet ./... && go test -test.short ./...



# Production Image
FROM base as builder
RUN go build -o openslides-search-service cmd/searchd/main.go


FROM alpine:3 as prod

ARG CONTEXT
ENV ${CONTEXT}=1

WORKDIR /

COPY entrypoint.sh ./
COPY meta/search.yml .
COPY meta/models.yml .
COPY --from=builder /app/openslides-search-service/openslides-search-service .

## External Information
LABEL org.opencontainers.image.title="OpenSlides Search Service"
LABEL org.opencontainers.image.description="The Search Service is a http endpoint where the clients can search for data within Openslides."
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.source="https://github.com/OpenSlides/openslides-search-service"

EXPOSE 9050

## Command
ENTRYPOINT ["./entrypoint.sh"]

CMD exec ./openslides-search-service