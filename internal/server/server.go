package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/alfrye/authorize/internal/router"
)

type Server struct {
	Engine *http.Server
	Router *router.Router
}

func New(port string) *Server {
	srv := &Server{
		Engine: &http.Server{
			Addr:         fmt.Sprintf("0.0.0.0:%s", port),
			WriteTimeout: time.Second * 15,
			ReadTimeout:  time.Second * 15,
			IdleTimeout:  time.Second * 15,
		},
		Router: router.New(),
	}

	return srv
}

func (s *Server) Listen() {

	fmt.Println("Starting Authorize API Server")
	log.Fatal(s.Engine.ListenAndServe())

}

func (s *Server) PopulateRoutes(routes []*router.API) {
	s.Router.LoadRoutes(routes)
	s.Engine.Handler = s.Router.Engine

}

func (s *Server) AttachOIDCRoutes(oidcRouter *mux.Router) {
	oidcRouter.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		s.Router.Engine.PathPrefix(route.GetName()).Handler(oidcRouter)
		return nil
	})

	muxRouter := mux.NewRouter()
	muxRouter.PathPrefix("/").Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		oidcRouter.ServeHTTP(w, r)
	}))

	s.Engine.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isOIDCPath(r.URL.Path) {
			oidcRouter.ServeHTTP(w, r)
		} else {
			s.Router.Engine.ServeHTTP(w, r)
		}
	})
}

func isOIDCPath(path string) bool {
	oidcPaths := []string{
		"/.well-known/openid-configuration",
		"/.well-known/jwks.json",
		"/authorize",
		"/token",
		"/introspect",
		"/userinfo",
	}
	for _, p := range oidcPaths {
		if path == p {
			return true
		}
	}
	return false
}
