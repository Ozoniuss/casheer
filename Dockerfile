# Build the application from source
FROM golang:1.21 AS build-stage

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

# hack to have this file in image always and pass e2e pipeline easily
RUN mkdir -m 666 /externaldeps
RUN touch /externaldeps/casheer.db

# Run the tests in the container (probably unnecessary if tests are run on
# every CICD run)
FROM build-stage AS run-test-stage
RUN go test -v ./...

# Deploy the application binary into a lean image
# Note that unlike the official go docker image guide, debian12 is required to
# avoid /lib/aarch64-linux-gnu/libc.so.6: version `GLIBC_2.3x' not found
# https://github.com/GoogleContainerTools/distroless/issues/1342
FROM gcr.io/distroless/base-debian12 AS build-release-stage

WORKDIR /

COPY --from=build-stage /casheer /casheer
COPY --from=build-stage --chown=nonroot:nonroot --chmod=666 /externaldeps/casheer.db /externaldeps/casheer.db

# These scripts are required to initialize the database if the db file
# doesn't exist. Note that creating the db file manually also requires
# running the scripts manually since the program only checks whether the
# db file exists or not.
#
# This step is included here and not in the docker-compose file since the 
# program will crash at startup if the database doesn't exist and it can't
# find these script files. Mounting the script files in docker-compose works,
# but this approach ships a "complete" image of casheer.
COPY ./scripts/sqlite /scripts/sqlite

USER nonroot:nonroot

ENTRYPOINT ["/casheer"]
CMD ["--server-address", "0.0.0.0", "--server-port", "8033"]