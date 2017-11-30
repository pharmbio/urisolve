package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"strings"

	"github.com/knakk/rdf"
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
		http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Just forward the raw RDF/XML from Blazegraph
	_, err = io.Copy(w, response.Body)
	if err != nil {
		http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
		return
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
		newTriples, err := h.runHdtQuery(uri + " ? ?")
		if err != nil {
			http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		triples = append(triples, newTriples...)
		newTriples, err = h.runHdtQuery("? ? " + uri)
		if err != nil {
			http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		triples = append(triples, newTriples...)
	}

	for _, triple := range triples {
		err := enc.Encode(triple)
		if err != nil {
			http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
	enc.Close()
}

func (h *UriResolverHandlerHdt) runHdtQuery(query string) ([]rdf.Triple, error) {
	var triples []rdf.Triple

	if !validQuery(query) {
		return nil, errors.New("Query contains invalid invalid characters")
	}

	Cmd := exec.Command("hdtSearch", "-q", query, h.HdtFilePath)
	hdtOut, err := Cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(hdtOut), "\n")
	for _, line := range lines {
		for _, l := range strings.Split(line, "\r") {
			if len(l) >= 4 && l[0:4] == "http" {
				triple, err := h.strToTriple(l)
				if err != nil {
					return nil, err
				}
				triples = append(triples, triple)
			}
		}
	}

	return triples, nil
}

func validQuery(query string) bool {
	charClass := `[A-Za-z0-9:\/\.\_\#\%]`
	validPattern := `(\?|` + charClass + `+) (\?|` + charClass + `+) (\?|` + charClass + `+)`
	validRegexp, err := regexp.Compile(validPattern)
	if err != nil {
		log.Fatal("Invalid regex: %s", validPattern)
	}
	return validRegexp.MatchString(query)
}

func (h *UriResolverHandlerHdt) strToTriple(line string) (rdf.Triple, error) {
	var triple rdf.Triple

	terms := strings.Split(line, " ")
	if len(terms) >= 3 {
		sRaw := terms[0]
		pRaw := terms[1]
		oRaw := terms[2]

		s, err := rdf.NewIRI(sRaw)
		if err != nil {
			return rdf.Triple{}, fmt.Errorf("Could not convert subject to IRI: %s (%s)", sRaw, err.Error())
		}

		p, err := rdf.NewIRI(pRaw)
		if err != nil {
			return rdf.Triple{}, fmt.Errorf("Could not convert predicate to IRI: %s (%s)", pRaw, err.Error())
		}

		if oRaw[0:1] == "h" {
			o, err := rdf.NewIRI(oRaw)
			if err != nil {
				return rdf.Triple{}, fmt.Errorf("Could not convert object to IRI: %s (%s)", oRaw, err.Error())
			}
			triple = rdf.Triple{s, p, o}
		} else if oRaw[0:1] == "\"" {
			o, err := rdf.NewLiteral(oRaw)
			if err != nil {
				return rdf.Triple{}, fmt.Errorf("Could not convert object to Literal: %s (%s)", oRaw, err.Error())
			}
			triple = rdf.Triple{s, p, o}
		}
	}

	return triple, nil
}
