package urlshort

import (
	"log"
	"net/http"

	"gopkg.in/yaml.v3"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	handler := func(w http.ResponseWriter, r *http.Request) {
		path := "/" + r.URL.Path[1:]
		if url, present := pathsToUrls[path]; present {
			http.Redirect(w, r, url, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
	return handler
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var pathUrls []pathURL

	err := yaml.Unmarshal(yml, &pathUrls)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		path := "/" + r.URL.Path[1:]

		for _, url := range pathUrls {
			if url.Path == path {
				http.Redirect(w, r, url.URL, http.StatusFound)
				return
			}
			fallback.ServeHTTP(w, r)
			return
		}
	}
	return handler, nil
}

type pathURL struct {
	Path string
	URL  string `yaml:"url"`
}
