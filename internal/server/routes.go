package server

import (
	//	"github.com/alfrye/authorize/internal/handlers/authorizeservice"
	"github.com/alfrye/authorize/internal/handlers/api"
	"github.com/alfrye/authorize/internal/router"
)

// AuthorizeServiceRoutes defines the routes
func (s *Server) AuthorizeServiceRoutes(handler api.AuthHandler) []*router.API {
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
							Handler: handler.RegisterUsers(),
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
							Handler: handler.Serve(),
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
							Handler: handler.Serve(),
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
							Handler: handler.Login(),
						},
					},
				},
			},
		},
	}
}
