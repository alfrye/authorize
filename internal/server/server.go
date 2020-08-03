package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alfrye/authorize/internal/router"
)

// Server represents a server instance
type Server struct {
	Engine *http.Server
	Router *router.Router
}

// New instaniates an new server instance
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

// Listen Starts the Web Server
func (s *Server) Listen() {

	fmt.Println("Starting Authorize API Server")
	log.Fatal(s.Engine.ListenAndServe())

}

func (s *Server) PopulateRoutes(routes []*router.API) {
	s.Router.LoadRoutes(routes)
	s.Engine.Handler = s.Router.Engine

}
