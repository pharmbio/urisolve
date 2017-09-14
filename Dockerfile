# Based on https://blog.golang.org/docker
# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:1.9

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/golang/pharmbio/urisolve

# Build the urisolve command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN go install github.com/golang/pharmbio/urisolve

# Run the urisolve command by default when the container starts.
CMD ["/go/bin/urisolve"]

# Document that the service listens on port 8888.
EXPOSE 8888
