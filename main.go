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
	host := flag.String("host", "localhost", "Hostname where to run this service")
	port := flag.String("port", "8888", "Port for this service")

	flag.Parse()
	if *endpoint == "" {
		log.Fatal("No SPARQL Endpoint URL provided. Use the -h flag to view options")
	}

	fmt.Println("Connecting to SPARQL Endpoint with URL:", *endpoint)

	// Main code
	fmt.Println("Starting to serve at: " + *host + ":" + *port + " ...")

	uriResHandler := &UriResolverHandler{*endpoint}

	http.Handle("/", uriResHandler)
	http.ListenAndServe(*host+":"+*port, nil)
}

// --------------------------------------------------------------------------------
// Uri Resolver Handler
// --------------------------------------------------------------------------------
type UriResolverHandler struct {
	SparqlEndpointUrl string
}

func (h *UriResolverHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Showing results from SPARQL endpoint at: %s\n\n", r.URL.Path[1:])

	repo, err := sparql.NewRepo(h.SparqlEndpointUrl,
		sparql.DigestAuth("", ""),
		sparql.Timeout(time.Millisecond*1500),
	)

	if err != nil {
		log.Fatal(err)
	}
	res, err := repo.Query("SELECT * WHERE { ?s ?p ?o } LIMIT 1")
	if err != nil {
		log.Println(err)
	}

	fmt.Fprintln(w, res)
}
