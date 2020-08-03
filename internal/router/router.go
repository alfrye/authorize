package router

import (
	"github.com/gorilla/mux"

	"net/http"
)

type (
	// Router handles incoming HTTP Request
	Router struct {
		Engine *mux.Router
		API    []*API
	}
	// API represents an API
	API struct {
		ServiceName string
		Version     string
		Subroutes   []*SubRoute
	}
	// SubRoute respresents the structure for sub routes
	SubRoute struct {
		Path      string
		Endpoints []*Endpoint
	}
	// Endpoint represents the endpoint for an API
	Endpoint struct {
		Method  string
		Path    string
		Queries []string
		Headers []string
		Handler http.HandlerFunc
	}
)

// New returns new router instance
func New() *Router {

	r := mux.NewRouter()

	router := &Router{
		Engine: r,
	}
	return router
}

func (r *Router) LoadRoutes(api []*API) {

	r.API = api
	for _, api := range r.API {

		// Creates a sub routes for service name amd version
		v := r.Engine.PathPrefix(api.ServiceName).PathPrefix(api.Version).Subrouter()

		for _, subroute := range api.Subroutes {
			s := v.PathPrefix(subroute.Path).Subrouter()
			for _, endpoint := range subroute.Endpoints {
				s.Methods(endpoint.Method).
					Headers(endpoint.Headers...).
					Queries(endpoint.Queries...).
					Path(endpoint.Path).
					Handler(endpoint.Handler)
			}
		}
	}
}
