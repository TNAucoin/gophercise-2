package main

import (
	"flag"
	"fmt"
	"github.com/tnaucoin/gophercise-2/urlshort"
	"log"
	"net/http"
	"os"
)

func main() {
	ymlPath := flag.String("ymlPath", "paths.yml", "A path to a file of yml URLs and short-paths")
	flag.Parse()
	yml, err := readYmlFile(*ymlPath)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yamlHandler, err := urlshort.YAMLHandler(yml, mapHandler)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", yamlHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

func readYmlFile(path string) ([]byte, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return f, nil
}
