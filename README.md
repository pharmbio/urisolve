# urisolve

A simple web server that resolving RDF URIs and returns RDF with any triples
connected to the URI in question.

## Installation

### Build from source

For now, building from source is the only available option.

1. [Install Go](https://golang.org/doc/install)
2. Then run this command:

   ```bash
   go get github.com/pharmbio/urisolve
   ```

## Usage

Given that you are running urisolve on a host that has the domain `example.org`
pointing to it, and you want to resolve URI:s starting with `example.org`,
which are available in a SPARQL endpoint with URL `endpoint-example.org`, you
could run it like this (sudo might be needed in order to bind to port 80, and
not required if you use a port above 1024):

```bash
sudo urisolve \
    -host example.org \
    -port 80 \
    -urihost http://example.org \
    -endpoint http://endpoint-example.org
```

To view the options available, run:

```bash
urisolve -h
```
