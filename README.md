# urisolve

A simple web server that resolves RDF URIs and returns RDF with any triples
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

The below examples assume, for demonstratory purposes, that you are running
urisolve on a host that has the domain `example.org` pointing to it, and you
want to resolve URI:s starting with `example.org`, which are available in your
data source (be that a SPARQL endpoint or an HDT file).

### With SPARQL endpoint as data source

Given that your preferred data source is a SPARQL endpoint with URL
`endpoint-example.org`, you could run it like this (sudo might be needed in
to bind to port 80, but typically not required if you use a port above 1024):

```bash
urisolve \
    -srctype sparql \
    -endpoint http://endpoint-example.org \
    -urihost http://example.org \
    -host example.org \
    -port 8080
```

### With HDT file as data source

If, instead of a SPARQL endpoint, you want to use an [(RDF) HDT](http://www.rdfhdt.org)
file as a data source (this ia an increasingly interesting option, as the
tooling around HDT matures), you can do it like this (not that this requires
the [C++ version of HDT tools](https://github.com/rdfhdt/hdt-cpp) installed):

```bash
urisolve \
    -srctype hdt \
    -hdtfile example_dataset.hdt \
    -urihost http://example.org \
    -host example.org \
    -port 8080
```

### More options

To view the options available, run:

```bash
urisolve -h
```
