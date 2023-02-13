package handler

import (
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
	return createMapHandler(pathsToUrls, fallback)
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

type mappedRequest struct {
	Path string
	Url  string
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml := make([]mappedRequest, 0)
	parseErr := yaml.Unmarshal(yml, &parsedYaml)

	if parseErr != nil {
		return nil, parseErr
	}

	pathsToUrls := make(map[string]string)

	for i := 0; i < len(parsedYaml); i++ {
		yamlEntry := parsedYaml[i]
		pathsToUrls[yamlEntry.Path] = yamlEntry.Url
	}

	return createMapHandler(pathsToUrls, fallback), nil
}

func createMapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		mappedUrl, hasUrl := pathsToUrls[path]

		if !hasUrl {
			fallback.ServeHTTP(w, r)
		} else {
			http.Redirect(w, r, mappedUrl, http.StatusMovedPermanently)
		}
	}
}
