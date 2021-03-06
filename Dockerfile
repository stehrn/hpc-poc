# Global so available across both (multi-stage) builds
ARG PACKAGE

# Use the offical golang image to create a binary.
# This is based on Debian and sets the GOPATH to /go.
# https://hub.docker.com/_/golang
FROM golang:1.15-buster as builder

ARG PACKAGE
ENV PACKAGE=${PACKAGE}
RUN echo "Running base build for package ${PACKAGE}"

# Create and change to the app directory.
WORKDIR /app

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
# Expecting to copy go.mod and if present go.sum.

COPY go.* ./
RUN go mod download

# Copy local code to the container image.
COPY . ./

# Build the binary.
RUN go build -mod=readonly -v -o server github.com/stehrn/hpc-poc/cmd/${PACKAGE}

# Use the official Debian slim image for a lean production container.
# https://hub.docker.com/_/debian
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM debian:buster-slim
RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

ARG PACKAGE
RUN echo "Running runtime build for package ${PACKAGE}"

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/server /app/server

# Run the service on container startup.
CMD ["/app/server"]