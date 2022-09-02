# Accept the Go version for the image to be set as a build argument.
# Default to Go 1.13
ARG GO_VERSION=1.13

# Dev stage
FROM golang:${GO_VERSION} AS dev_builder
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN apt-get update -qq
RUN apt-get clean && rm -rf /var/lib/apt/lists/*
WORKDIR /go/src/github.com/mrzacarias/example-app
COPY . .

# First stage: build the executable.
FROM golang:${GO_VERSION}-alpine AS builder

# Create the user and group files that will be used in the running container to
# run the process as an unprivileged user.
RUN mkdir /user && \
  echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
  echo 'nobody:x:65534:' > /user/group

# Install the Certificate-Authority certificates for the app to be able to make
# calls to HTTPS endpoints.
# Git is required for fetching the dependencies.
RUN apk add --no-cache ca-certificates git

# Set the working directory to $GOPATH
WORKDIR /go/src/github.com/mrzacarias/example-app

# Import the code from the context.
COPY ./ ./

# Build the executable to `/app`. Mark the build as statically linked.
RUN CGO_ENABLED=0 go build \
  -installsuffix 'static' \
  -o /app ./cmd/app

# Final stage: the running container.
FROM scratch AS final

# Import the user and group files from the first stage.
COPY --from=builder /user/group /user/passwd /etc/

# Import the Certificate-Authority certificates for enabling HTTPS.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Import the compiled executable from the first stage.
COPY --from=builder /app /app

# Declare the port on which the webserver will be exposed.
# As we're going to run the executable as an unprivileged user, we can't bind
# to ports below 1024.
EXPOSE 8080 8081

# Perform any further action as an unprivileged user.
USER 65534:65534

# Run the compiled binary.
ENTRYPOINT ["/app"]
