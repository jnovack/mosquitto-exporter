# build stage
FROM golang:1.14-alpine AS build
WORKDIR /go/src/app

RUN apk add --no-cache git ca-certificates

COPY go.mod .
COPY go.sum .
RUN go mod download

ARG APPLICATION="go-executable"
ARG BUILD_RFC3339="1970-01-01T00:00:00Z"
ARG COMMIT="local"
ARG VERSION="dirty"
ARG GO_LDFLAGS="-w -s \
        -X github.com/jnovack/go-version.Application=${APPLICATION} \
        -X github.com/jnovack/go-version.BuildDate=${BUILD_RFC3339} \
        -X github.com/jnovack/go-version.Revision=${COMMIT} \
        -X github.com/jnovack/go-version.Version=${VERSION} \
        "

# Build
COPY . .
RUN go build -ldflags "${GO_LDFLAGS}" -o /go/bin/${APPLICATION} .

###############################################################################
# final stage
FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ARG BUILD_RFC3339="1970-01-01T00:00:00Z"
ARG COMMIT="local"
ARG VERSION="dirty"

LABEL org.opencontainers.image.ref.name="jnovack/mosquitto-exporter" \
      org.opencontainers.image.created=$BUILD_RFC3339 \
      org.opencontainers.image.authors="Justin J. Novack <jnovack@gmail.com>" \
      org.opencontainers.image.documentation="https://github.com/jnovack/mosquitto-exporter/README.md" \
      org.opencontainers.image.description="Minimalist mosquitto-exporter for single host deployments." \
      org.opencontainers.image.licenses="MIT" \
      org.opencontainers.image.source="https://github.com/jnovack/mosquitto-exporter" \
      org.opencontainers.image.revision=$COMMIT \
      org.opencontainers.image.version=$VERSION \
      org.opencontainers.image.url="https://hub.docker.com/r/jnovack/mosquitto-exporter/"

COPY --from=build /go/bin/mosquitto-exporter /mosquitto-exporter

EXPOSE 9344

ENTRYPOINT ["/mosquitto-exporter"]
