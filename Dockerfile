# Build the application from source
FROM golang:1.24 AS build-stage

WORKDIR /app

# For the moment, unnecessary.
COPY go.mod go.sum .
RUN go mod download

COPY ./cmd/. cmd
COPY ./currency currency
COPY ./internal/. internal
COPY ./pkg/. pkg

# CGO_ENABLED must be 1 because go-sqlite3 is a CGO enabled package: 
# https://github.com/mattn/go-sqlite3#installation
RUN CGO_ENABLED=1 GOOS=linux go build -o /casheer ./cmd/

# Deploy the application binary into a lean image
# Note that unlike the official go docker image guide, debian12 is required to
# avoid /lib/aarch64-linux-gnu/libc.so.6: version `GLIBC_2.3x' not found
# https://github.com/GoogleContainerTools/distroless/issues/1342
FROM gcr.io/distroless/base-debian12 AS build-release-stage

WORKDIR /

COPY --from=build-stage /casheer /casheer

USER nonroot:nonroot

ENTRYPOINT ["/casheer"]
CMD ["--server-address", "0.0.0.0", "--server-port", "8033"]