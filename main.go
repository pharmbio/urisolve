package main

import (
	"flag"
	"fmt"
	"github.com/knakk/rdf"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

func main() {
	// Set up flags
	srcType := flag.String("srctype", "", "Type of data source. Can be one of: sparql, hdt")
	urihost := flag.String("urihost", "", "Hostname for which to resolve URIs (without trailing slash)")
	endpoint := flag.String("endpoint", "", "URL to a SPARQL 1.1 endpoint")
	host := flag.String("host", "localhost", "Hostname where to run this service (without trailing slash)")
	port := flag.String("port", "8080", "Port where this service should be exposed")
	hdtFilePath := flag.String("hdtfile", "", "A (relative or full) path to an .hdt file")

	// Parse flags
	flag.Parse()

	// Handle flag errors
	if *srcType == "sparql" {
		if *endpoint == "" {
			log.Fatal("No SPARQL Endpoint URL provided. Use the -h flag to view options")
		}
	} else if *srcType == "hdt" {
		if *hdtFilePath == "" {
			log.Fatal("No HDT file path specified! You have to specify a path to a .hdt file using the -hdtfile flag. Use -h to view options")
		}
	} else {
		log.Fatal("Invalid source type specified. You have to use the -srctype flag to specify either 'sparql' or 'hdt'. Use -h to view options")
	}

	if *urihost == "" {
		log.Fatal("No urihost provided. Use the -h flag to view options")
	}

	// Execute the relevant HTTP handler, based on the source type selected
	if *srcType == "sparql" {
		// Print some output to the console
		fmt.Println("Connecting to SPARQL Endpoint with URL:", *endpoint)
		fmt.Println("Starting to serve at: " + *host + ":" + *port + " ...")

		// Start handling requests
		uriResHandlerSparql := &UriResolverHandlerSparql{*urihost, *endpoint}
		http.Handle("/", uriResHandlerSparql)
	} else if *srcType == "hdt" {
		// Print some output to the console
		fmt.Println("Using the following HDT for querying: ", *hdtFilePath)
		fmt.Println("Starting to serve at: " + *host + ":" + *port + " ...")

		// Start handling requests
		uriResHandlerHdt := &UriResolverHandlerHdt{*urihost, *hdtFilePath}
		http.Handle("/", uriResHandlerHdt)
	}

	// Start serving requests
	err := http.ListenAndServe(*host+":"+*port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

// UriResolverHandlerSparql handles RDF URI:s and writes out RDF with any triples
// connected to the URI in question, to w, based on information in a SPARQL
// endpoint as indicated with the SparqlEndpointUrl field, which has to be set
// upon creating a new UriResolverHandlerSparql.
type UriResolverHandlerSparql struct {
	UriHost           string
	SparqlEndpointUrl string
}

func (h *UriResolverHandlerSparql) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	uri := h.UriHost + "/" + r.URL.Path[1:]

	sparqlQuery := `query=DESCRIBE <` + uri + `>`

	fmt.Println("Querying " + h.SparqlEndpointUrl + " with the following parameters:")
	fmt.Println(sparqlQuery)

	reader := strings.NewReader(sparqlQuery)
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
	_, err = io.Copy(w, response.Body)
	if err != nil {
		log.Fatal(err)
	}
}

// UriResolverHandlerHdt handles RDF URI:s and writes out RDF with any triples
// connected to the URI in question, to w, based on information in a (RDF)HDT
// dataset file. You can find more info about hDT at http://www.rdfhdt.org
type UriResolverHandlerHdt struct {
	UriHost     string
	HdtFilePath string
}

func (h *UriResolverHandlerHdt) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	enc := rdf.NewTripleEncoder(w, rdf.NTriples)

	var triples []rdf.Triple

	path := r.URL.Path[1:]
	if path != "favicon.ico" {
		uri := h.UriHost + "/" + r.URL.Path[1:]
		triples = append(triples, h.runHdtQuery(uri+" ? ?")...)
		triples = append(triples, h.runHdtQuery("? ? "+uri)...)
	}

	for _, triple := range triples {
		err := enc.Encode(triple)
		if err != nil {
			log.Fatalf("Could not parse triple: %v\n%v", triple, err)
		}
	}
	enc.Close()
}

func (h *UriResolverHandlerHdt) runHdtQuery(query string) []rdf.Triple {
	var triples []rdf.Triple

	Cmd := exec.Command("hdtSearch", "-q", query, h.HdtFilePath)
	hdtOut, err := Cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(hdtOut), "\n")
	for _, line := range lines {
		for _, l := range strings.Split(line, "\r") {
			if len(l) >= 4 && l[0:4] == "http" {
				triples = append(triples, h.strToTriple(l))
			}
		}
	}

	return triples
}

func (h *UriResolverHandlerHdt) strToTriple(line string) rdf.Triple {
	var triple rdf.Triple

	terms := strings.Split(line, " ")
	if len(terms) >= 3 {
		sRaw := terms[0]
		pRaw := terms[1]
		oRaw := terms[2]

		s, err := rdf.NewIRI(sRaw)
		if err != nil {
			log.Fatalf("Could not convert subject to IRI: %s\n", sRaw)
		}

		p, err := rdf.NewIRI(pRaw)
		if err != nil {
			log.Fatalf("Could not convert predicate to IRI: %s\n", pRaw)
		}

		if oRaw[0:1] == "h" {
			o, err := rdf.NewIRI(oRaw)
			if err != nil {
				log.Fatalf("Could not convert object to IRI: %s\n", oRaw)
			}
			triple = rdf.Triple{s, p, o}
		} else if oRaw[0:1] == "\"" {
			o, err := rdf.NewLiteral(oRaw)
			if err != nil {
				log.Fatalf("Could not convert object to Literal: %s\n", oRaw)
			}
			triple = rdf.Triple{s, p, o}
		}
	}

	return triple
}
