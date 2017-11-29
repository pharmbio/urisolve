#!/bin/bash
echo "After the service started below, open a browser and surf in on http://localhost:8080/entity/Q3 ..."
docker run \
    -it \
    -p 8080:8080 \
    --name urisolve-test-sparql \
    --hostname urisolve \
    --rm \
    farmbio/urisolve \
    urisolve \
    -srctype sparql \
    -endpoint https://query.wikidata.org/sparql \
    -urihost http://www.wikidata.org \
    -host urisolve \
    -port 8080
