# Base this image on https://github.com/rdfhdt/hdt-docker ... so that we get
# the HDT C++ tools installed by default
FROM rfdhdt/hdt-cpp

# Set up environment variables
ENV GOROOT=/opt/go
ENV GOPATH=/usr/go
ENV PATH="/usr/go/bin:/opt/go/bin:${PATH}"
ENV HDTFILE=/usr/go/src/github.com/pharmbio/urisolve/example_data.hdt

# We need rsync to be able to move data into pod
RUN apt-get install -y rsync

# Install a more up to date version of Go (that supports vendoring etc)
RUN cd /opt && curl -O https://storage.googleapis.com/golang/go1.9.linux-amd64.tar.gz && tar -zxf go1.9.linux-amd64.tar.gz

# Copy the local package files to the container's workspace.
ADD . /usr/go/src/github.com/pharmbio/urisolve

# Build the urisolve command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN go install github.com/pharmbio/urisolve

# We must not run as root, because of OpenShift policy
USER 1001

# Run the urisolve command by default when the container starts.
CMD urisolve -srctype hdt -hdtfile $HDTFILE -urihost http://rdf.pharmb.io -host $HOSTNAME -port 8080

# Document that the service listens on port 8080.
EXPOSE 8080
