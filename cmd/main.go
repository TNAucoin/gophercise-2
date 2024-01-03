package main

import (
	"flag"
	"fmt"
	"github.com/tnaucoin/gophercise-2/urlshort"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	filePath := flag.String("filePath", "paths.json", "A path to a file of yml URLs and short-paths")
	fileType := flag.String("fileType", "json", "Type of file data [json,yaml,yml]")
	flag.Parse()
	fd, err := readData(*filePath)
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

	var handler http.HandlerFunc
	switch strings.ToLower(strings.TrimSpace(*fileType)) {
	case "json":
		handler, err = urlshort.JSONHandler(fd, mapHandler)
		if err != nil {
			panic(err)
		}
	case "yml":
	case "yaml":
		// Build the YAMLHandler using the mapHandler as the
		// fallback
		handler, err = urlshort.YAMLHandler(fd, mapHandler)
		if err != nil {
			panic(err)
		}
	default:
		log.Println("invalid file type, using default handler")
		handler = mapHandler
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", handler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

func readData(path string) ([]byte, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return f, nil
}
