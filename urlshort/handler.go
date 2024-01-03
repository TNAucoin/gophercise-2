package urlshort

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
	"net/http"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	//	TODO: Implement this...
	return func(w http.ResponseWriter, r *http.Request) {
		url, ok := pathsToUrls[r.URL.Path]
		if ok {
			http.Redirect(w, r, url, http.StatusMovedPermanently)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
type PathUrlData struct {
	Path string `yaml:"path" json:"path"`
	Url  string `yaml:"url" json:"url"`
}

func unmarshalYamlData(yml []byte) ([]PathUrlData, error) {
	yd := []PathUrlData{}
	if err := yaml.Unmarshal(yml, &yd); err != nil {
		return nil, err
	}
	return yd, nil
}

func unmarshalJSONData(data []byte) ([]PathUrlData, error) {
	js := []PathUrlData{}
	if err := json.Unmarshal(data, &js); err != nil {
		return nil, err
	}
	return js, nil
}

func createPathUrlMap(yamlData []PathUrlData) map[string]string {
	ym := make(map[string]string)
	for _, v := range yamlData {
		ym[v.Path] = v.Url
	}
	return ym
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	// TODO: Implement this...
	yd, err := unmarshalYamlData(yml)
	if err != nil {
		return nil, err
	}
	ym := createPathUrlMap(yd)
	return MapHandler(ym, fallback), nil
}

func JSONHandler(json []byte, fallback http.Handler) (http.HandlerFunc, error) {
	js, err := unmarshalJSONData(json)
	if err != nil {
		return nil, err
	}
	jm := createPathUrlMap(js)
	return func(w http.ResponseWriter, r *http.Request) {
		p, ok := jm[r.URL.Path]
		if ok {
			http.Redirect(w, r, p, http.StatusMovedPermanently)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}, nil
}
