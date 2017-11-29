#!/bin/bash
export HDTFILE=/usr/go/src/github.com/pharmbio/urisolve/example_data.hdt
echo "After the service started below, open a browser and surf in on http://localhost:8080/cplogd/Compound1 ..."
docker run \
    -it \
    -p 8080:8080 \
    --name urisolve-test-hdt \
    --hostname urisolve \
    --rm \
    farmbio/urisolve

# Full command, if you want to customize some parameters to the urisolve binary
#docker run \
#    -it \
#    -p 8080:8080 \
#    --name urisolve-test-hdt \
#    --hostname urisolve \
#    --rm farmbio/urisolve \
#    urisolve \
#    -srctype hdt \
#    -hdtfile $HDTFILE \
#    -urihost http://rdf.pharmb.io \
#    -host urisolve \
#    -port 8080
