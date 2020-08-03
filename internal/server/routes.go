package server

import (
	//	"github.com/alfrye/authorize/internal/handlers/authorizeservice"
	"github.com/alfrye/authorize/internal/handlers/authorizeservice"
	"github.com/alfrye/authorize/internal/router"
)

// AuthorizeServiceRoutes defines the routes
func (s *Server) AuthorizeServiceRoutes() []*router.API {
	return []*router.API{
		{
			ServiceName: "/authorize",
			Version:     "/v1",
			Subroutes: []*router.SubRoute{
				{
					Path: "/register",
					Endpoints: []*router.Endpoint{
						{
							Method:  "POST",
							Path:    "",
							Handler: authorizeservice.RegisterUsers(),
						},
					},
				},
			},
		},
		{
			ServiceName: "/authorize",
			Version:     "/v1",
			Subroutes: []*router.SubRoute{
				{
					Path: "/user",
					Endpoints: []*router.Endpoint{
						{
							Method:  "GET",
							Path:    "",
							Handler: authorizeservice.Serve(),
						},
					},
				},
			},
		},

		{
			ServiceName: "/authorize",
			Version:     "/v1",
			Subroutes: []*router.SubRoute{
				{
					Path: "/user",
					Endpoints: []*router.Endpoint{
						{
							Method:  "POST",
							Path:    "",
							Handler: authorizeservice.Serve(),
						},
					},
				},
			},
		},
		{
			ServiceName: "/authorize",
			Version:     "/v1",
			Subroutes: []*router.SubRoute{
				{
					Path: "/login",
					Endpoints: []*router.Endpoint{
						{
							Method:  "POST",
							Path:    "",
							Handler: authorizeservice.Login(),
						},
					},
				},
			},
		},
	}
}
