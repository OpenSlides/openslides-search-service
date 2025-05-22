ARG CONTEXT=prod
ARG GO_IMAGE_VERSION=1.24.3
ARG ALPINE_VERSION=3

FROM golang:${GO_IMAGE_VERSION}-alpine as base

ARG CONTEXT
ARG GO_IMAGE_VERSION

WORKDIR /root/openslides-search-service

## Installs
RUN apk add git

COPY go.mod go.sum ./
RUN go mod download

COPY cmd cmd
COPY pkg pkg

LABEL org.opencontainers.image.title="OpenSlides Search Service"
LABEL org.opencontainers.image.description="The Search Service is a http endpoint where the clients can search for data within Openslides."
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.source="https://github.com/OpenSlides/openslides-search-service"

EXPOSE 9050

# Development Image

FROM base as dev

RUN ["go", "install", "github.com/githubnemo/CompileDaemon@latest"]

WORKDIR /root
COPY entrypoint.sh ./
COPY meta/search.yml .
COPY meta/models.yml .
ENTRYPOINT ["./entrypoint.sh"]

CMD CompileDaemon -log-prefix=false -build="go build -o search-service ./openslides-search-service/cmd/searchd/main.go" -command="./search-service"

# Testing Image

FROM base as tests

RUN apk add build-base

CMD go vet ./... && go test -test.short ./...


# Production Image

FROM base as builder
RUN go build -o openslides-search-service cmd/searchd/main.go


FROM scratch as prod

COPY entrypoint.sh ./
COPY meta/search.yml .
COPY meta/models.yml .
COPY --from=builder /root/openslides-search-service/openslides-search-service .

EXPOSE 9050
ENTRYPOINT ["./entrypoint.sh"]
CMD exec ./openslides-search-service