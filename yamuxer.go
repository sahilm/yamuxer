package yamuxer

import (
	"context"
	"net/http"
	"regexp"
)

type Route struct {
	Pattern *regexp.Regexp
	Handler func(w http.ResponseWriter, r *http.Request)
}

type Mux struct {
	routes []*Route
}

func New(routes []*Route) *Mux {
	return &Mux{routes}
}

func (mux Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range mux.routes {
		matches := route.Pattern.FindStringSubmatch(r.URL.Path)
		if len(matches) > 0 {
			ctx := context.WithValue(r.Context(), ContextKey("matches"), matches[1:])
			route.Handler(w, r.WithContext(ctx))
			return
		}
	}
	http.NotFound(w, r)
}

type ContextKey string
