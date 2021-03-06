# build stage
FROM golang:1.15-alpine AS build
WORKDIR /go/src/app

# Create appuser so application does not run as root in final container
# See https://stackoverflow.com/a/55757473/12429735
ENV USER=appuser
ENV UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

# Add ca-certificates in the event final application needs ssl/tls
RUN apk add --no-cache git ca-certificates

# Make use of build caching by first downloading all dependancies
COPY go.mod .
COPY go.sum .
RUN go mod download

# Set up docker build arguments
ARG APPLICATION="myapp"
ARG BUILD_RFC3339="1970-01-01T00:00:00Z"
ARG REVISION="local"
ARG VERSION="dirty"
ARG GO_LDFLAGS="-w -s \
    -X github.com/jnovack/release.Application=${APPLICATION} \
    -X github.com/jnovack/release.BuildDate=${BUILD_RFC3339} \
    -X github.com/jnovack/release.Revision=${REVISION} \
    -X github.com/jnovack/release.Version=${VERSION} \
    -extldflags '-static'"

# Build the main application
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "${GO_LDFLAGS}" -o /go/bin/${APPLICATION} cmd/${APPLICATION}/main.go

###############################################################################
# final stage
FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/group /etc/group
USER appuser:appuser

ARG APPLICATION="myapp"
ARG BUILD_RFC3339="1970-01-01T00:00:00Z"
ARG REVISION="local"
ARG DESCRIPTION="no description"
ARG PACKAGE="user/repo"
ARG VERSION="dirty"

LABEL org.opencontainers.image.ref.name="${PACKAGE}" \
    org.opencontainers.image.created=$BUILD_RFC3339 \
    org.opencontainers.image.authors="Justin J. Novack <jnovack@gmail.com>" \
    org.opencontainers.image.documentation="https://github.com/${PACKAGE}/README.md" \
    org.opencontainers.image.description="${DESCRIPTION}" \
    org.opencontainers.image.licenses="GPLv3" \
    org.opencontainers.image.source="https://github.com/${PACKAGE}" \
    org.opencontainers.image.revision=$REVISION \
    org.opencontainers.image.version=$VERSION \
    org.opencontainers.image.url="https://hub.docker.com/r/${PACKAGE}/"

COPY --from=build /go/bin/${APPLICATION} /app

EXPOSE 9344

ENTRYPOINT ["/app"]
