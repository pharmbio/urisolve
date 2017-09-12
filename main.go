package main

import (
	"flag"
	"fmt"
	"github.com/knakk/rdf"
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
	Namespace         string
	SparqlEndpointUrl string
}

func (h *UriResolverHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	uri := h.Namespace + r.URL.Path[1:]

	// Connect to SPARQL Endpoint
	repo, err := sparql.NewRepo(h.SparqlEndpointUrl,
		sparql.DigestAuth("", ""),
		sparql.Timeout(time.Millisecond*1500),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Run SPARQL query
	results, err := repo.Query("DESCRIBE <" + uri + ">")
	if err != nil {
		log.Println(err)
	}

	for _, sol := range results.Solutions() {
		subj, err := rdf.NewIRI(sol["subject"].String())
		if err != nil {
			log.Fatal(err)
		}

		pred, err := rdf.NewIRI(sol["predicate"].String())
		if err != nil {
			log.Fatal(err)
		}

		var obj rdf.Object
		switch sol["object"].Type() {
		case rdf.TermBlank:
			obj, err = rdf.NewBlank(sol["object"].String())
			if err != nil {
				log.Fatal(err)
			}
		case rdf.TermIRI:
			obj, err = rdf.NewIRI(sol["object"].String())
			if err != nil {
				log.Fatal(err)
			}
		case rdf.TermLiteral:
			obj, err = rdf.NewLiteral(sol["object"].String())
			if err != nil {
				log.Fatal(err)
			}
		}

		triple := rdf.Triple{subj, pred, obj}
		fmt.Fprint(w, triple.Serialize(rdf.Turtle))
	}
}
