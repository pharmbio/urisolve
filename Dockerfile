FROM rfdhdt/hdt-cpp

ENV GOROOT /opt/go
ENV GOPATH /usr/go
ENV HDTFILE example_data.hdt

RUN cd /opt && curl -O https://storage.googleapis.com/golang/go1.9.linux-amd64.tar.gz && tar -zxvf go1.9.linux-amd64.tar.gz

# Copy the local package files to the container's workspace.
ADD . /usr/go/src/github.com/pharmbio/urisolve

# Build the urisolve command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN go install github.com/pharmbio/urisolve

# Run the urisolve command by default when the container starts.
CMD /opt/go/bin/urisolve -srctype hdt -hdtfile $HDTFILE -urihost http://example.org

# Document that the service listens on port 8080.
EXPOSE 8080
