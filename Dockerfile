# Based on https://blog.golang.org/docker
# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# These are supposed to be overriden when running
ENV HDTFILE example_data.hdt

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/pharmbio/urisolve

# Build the urisolve command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN go get github.com/knakk/rdf
RUN go install github.com/pharmbio/urisolve

# Run the urisolve command by default when the container starts.
CMD ["/go/bin/urisolve", "-srctype", "hdt", "-hdtfile", "$HDTFILE", "-urihost", "http://example.org"]

# Document that the service listens on port 8080.
EXPOSE 8080
