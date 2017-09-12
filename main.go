package main

import (
	"flag"
	"fmt"
	"github.com/knakk/sparql"
	"log"
	"net/http"
	"time"
)

func main() {
	// Flag handling
	endpoint := flag.String("endpoint", "", "Provide an URL to a SPARQL 1.1 endpoint")
	namespace := flag.String("namespace", "", "Specify the URI namespace for which to resolve URIs")
	host := flag.String("host", "localhost", "Hostname where to run this service")
	port := flag.String("port", "8888", "Port for this service")
	flag.Parse()
	if *endpoint == "" {
		log.Fatal("No SPARQL Endpoint URL provided. Use the -h flag to view options")
	} else if *namespace == "" {
		log.Fatal("No namespace provided. Use the -h flag to view options")
	}

	// Some status info
	fmt.Println("Connecting to SPARQL Endpoint with URL:", *endpoint)
	fmt.Println("Starting to serve at: " + *host + ":" + *port + " ...")

	// Start handling requests
	uriResHandler := &UriResolverHandler{*namespace, *endpoint}
	http.Handle("/", uriResHandler)
	http.ListenAndServe(*host+":"+*port, nil)
}

// UriResolverHandler handles RDF URI:s and writes out RDF with any triples
// connected to the URI in question, to w, based on information in a SPARQL
// endpoint as indicated with the SparqlEndpointUrl field, which has to be set
// upon creating a new UriResolverHandler.
type UriResolverHandler struct {
	Namespace         string
	SparqlEndpointUrl string
}

func (h *UriResolverHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Some debugging output
	fmt.Fprintf(w, "Showing results from SPARQL endpoint at: %s\n\n", r.URL.Path[1:])

	// Connect to SPARQL Endpoint
	repo, err := sparql.NewRepo(h.SparqlEndpointUrl,
		sparql.DigestAuth("", ""),
		sparql.Timeout(time.Millisecond*1500),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Run SPARQL query
	uri := h.Namespace + r.URL.Path[1:]
	res, err := repo.Query("SELECT * WHERE { <" + uri + "> ?p ?o } LIMIT 1")
	if err != nil {
		log.Println(err)
	}

	// Print out results
	fmt.Fprintln(w, res)
}
