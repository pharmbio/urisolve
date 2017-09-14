package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func main() {
	// Flag handling
	endpoint := flag.String("endpoint", "", "Provide an URL to a SPARQL 1.1 endpoint")
	urihost := flag.String("urihost", "", "Specify the URI hostname for which to resolve URIs (without trailing slash)")
	host := flag.String("host", "localhost", "Hostname where to run this service (without trailing slash)")
	port := flag.String("port", "8888", "Port for this service")
	flag.Parse()
	if *endpoint == "" {
		log.Fatal("No SPARQL Endpoint URL provided. Use the -h flag to view options")
	} else if *urihost == "" {
		log.Fatal("No urihost provided. Use the -h flag to view options")
	}

	// Print some output to the console
	fmt.Println("Connecting to SPARQL Endpoint with URL:", *endpoint)
	fmt.Println("Starting to serve at: " + *host + ":" + *port + " ...")

	// Start handling requests
	uriResHandler := &UriResolverHandler{*urihost, *endpoint}
	http.Handle("/", uriResHandler)
	err := http.ListenAndServe(*host+":"+*port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

// UriResolverHandler handles RDF URI:s and writes out RDF with any triples
// connected to the URI in question, to w, based on information in a SPARQL
// endpoint as indicated with the SparqlEndpointUrl field, which has to be set
// upon creating a new UriResolverHandler.
type UriResolverHandler struct {
	UriHost           string
	SparqlEndpointUrl string
}

func (h *UriResolverHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	uri := h.UriHost + "/" + r.URL.Path[1:]

	queryString := `query=DESCRIBE <` + uri + `>`

	fmt.Println("Querying " + h.SparqlEndpointUrl + " with the following parameters:")
	fmt.Println(queryString)

	reader := strings.NewReader(queryString)
	request, err := http.NewRequest("POST", h.SparqlEndpointUrl, reader)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	// Just forward the raw RDF/XML from Blazegraph
	io.Copy(w, response.Body)
}
