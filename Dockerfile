ARG CONTEXT=prod
ARG ALPINE_VERSION=3

FROM golang:1.24.3-alpine as base

## Setup
ARG CONTEXT
WORKDIR /root/openslides-search-service
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

## Command
COPY ./dev/command.sh ./
RUN chmod +x command.sh
CMD ["./command.sh"]


# Development Image
FROM base as dev

RUN ["go", "install", "github.com/githubnemo/CompileDaemon@latest"]

WORKDIR /root
COPY entrypoint.sh ./
COPY meta/search.yml .
COPY meta/models.yml .

## Command
ENTRYPOINT ["./entrypoint.sh"]



# Testing Image
FROM base as tests

RUN apk add build-base --no-cache





# Production Image
FROM base as builder
RUN go build -o openslides-search-service cmd/searchd/main.go


FROM alpine:3 as prod

ARG CONTEXT

WORKDIR /

COPY entrypoint.sh ./
COPY meta/search.yml .
COPY meta/models.yml .
COPY --from=builder /root/openslides-search-service/openslides-search-service .

## External Information
LABEL org.opencontainers.image.title="OpenSlides Search Service"
LABEL org.opencontainers.image.description="The Search Service is a http endpoint where the clients can search for data within Openslides."
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.source="https://github.com/OpenSlides/openslides-search-service"

EXPOSE 9050

## Command
ENV ${CONTEXT}=1
COPY ./dev/command.sh ./
RUN chmod +x command.sh
CMD ["./command.sh"]
ENTRYPOINT ["./entrypoint.sh"]