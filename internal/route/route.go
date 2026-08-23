package route

import (
	"net/http"
	"strings"
)

type Match struct {
	Method, Path string
	Handler      http.Handler
}
type Table struct{ matches []Match }

func (t *Table) Add(method, path string, handler http.Handler) {
	t.matches = append(t.matches, Match{Method: method, Path: path, Handler: handler})
}
func (t *Table) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, match := range t.matches {
		if match.Method == r.Method && matchPath(match.Path, r.URL.Path) {
			match.Handler.ServeHTTP(w, r)
			return
		}
	}
	http.NotFound(w, r)
}
func matchPath(pattern, path string) bool {
	patternParts := strings.Split(strings.Trim(pattern, "/"), "/")
	pathParts := strings.Split(strings.Trim(path, "/"), "/")
	if len(patternParts) != len(pathParts) {
		return false
	}
	for i, part := range patternParts {
		if strings.HasPrefix(part, ":") {
			continue
		}
		if part != pathParts[i] {
			return false
		}
	}
	return true
}
