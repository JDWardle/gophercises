package urlshort

import (
	"encoding/json"
	"net/http"

	"gopkg.in/yaml.v2"
)

type pathToURL struct {
	Path string `yaml:"path" json:"path"`
	URL  string `yaml:"url" json:"url"`
}

func pathsToURLMap(paths []pathToURL) map[string]string {
	pathsToUrls := map[string]string{}
	for _, path := range paths {
		pathsToUrls[path.Path] = path.URL
	}
	return pathsToUrls
}

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.EscapedPath()

		if _, ok := pathsToUrls[path]; !ok {
			fallback.ServeHTTP(w, r)
			return
		}

		http.Redirect(w, r, pathsToUrls[path], http.StatusMovedPermanently)
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
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	paths := []pathToURL{}

	if err := yaml.Unmarshal(yml, &paths); err != nil {
		return nil, err
	}

	return MapHandler(pathsToURLMap(paths), fallback), nil
}

// JSONHandler parses a JSON string for URL paths that it will attempt
// to redirect to based on the URL path of the request.
func JSONHandler(jsonData []byte, fallback http.Handler) (http.HandlerFunc, error) {
	paths := []pathToURL{}

	if err := json.Unmarshal(jsonData, &paths); err != nil {
		return nil, err
	}

	return MapHandler(pathsToURLMap(paths), fallback), nil
}
