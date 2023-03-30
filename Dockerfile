FROM golang:1.20.2-alpine as base
WORKDIR /root/

RUN apk add git

COPY go.mod go.sum ./
RUN go mod download

COPY cmd cmd
COPY pkg pkg

# Build service in seperate stage.
FROM base as builder
RUN go build -o openslides-search-service cmd/searchd/main.go
RUN go build -o openslides-search-generate-filter cmd/generate-filter/main.go


# Test build.
FROM base as testing

RUN apk add build-base

CMD go vet ./... && go test -test.short ./...


# Development build.
FROM base as development

COPY --from=builder /root/openslides-search-generate-filter /root/openslides-search-generate-filter
RUN ./openslides-search-generate-filter --output search.yml 
RUN ["go", "install", "github.com/githubnemo/CompileDaemon@latest"]
EXPOSE 9012

RUN wget https://github.com/OpenSlides/openslides-backend/raw/main/global/meta/models.yml
CMD CompileDaemon -log-prefix=false -build="go build -o openslides-search-service cmd/searchd/main.go" -command="./openslides-search-service"


# Productive build
FROM scratch

LABEL org.opencontainers.image.title="OpenSlides Search Service"
LABEL org.opencontainers.image.description="The Search Service is a http endpoint where the clients can search for data within Openslides."
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.source="https://github.com/OpenSlides/openslides-search-service"

COPY --from=builder /root/openslides-search-service .
EXPOSE 9012
ENTRYPOINT ["/openslides-search-service"]
