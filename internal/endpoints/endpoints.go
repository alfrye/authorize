package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/alfrye/authorize/internal/models"
)

type (
	// AuthService defines the service interface for go-kit
	AuthService interface {
		Login(ctx context.Context, user models.Users) (string, error)
		Register(ctx context.Context, user models.Users) error
		GoogleCallback(ctx context.Context, code, state string) (models.Users, error)
		Serve(ctx context.Context, token string) (models.Users, error)
	}

	// Endpoints holds all the endpoints for the service
	Endpoints struct {
		LoginEndpoint        endpoint.Endpoint
		RegisterEndpoint      endpoint.Endpoint
		GoogleCallbackEndpoint endpoint.Endpoint
		ServeEndpoint        endpoint.Endpoint
	}

	// LoginRequest represents the login request
	LoginRequest struct {
		User models.Users `json:"user"`
	}

	// LoginResponse represents the login response
	LoginResponse struct {
		RedirectURL string `json:"redirect_url"`
		Err         string `json:"err,omitempty"`
	}

	// RegisterRequest represents the register request
	RegisterRequest struct {
		User models.Users `json:"user"`
	}

	// RegisterResponse represents the register response
	RegisterResponse struct {
		Message string `json:"message"`
		Err     string `json:"err,omitempty"`
	}

	// GoogleCallbackRequest represents the Google OAuth callback request
	GoogleCallbackRequest struct {
		Code  string `json:"code"`
		State string `json:"state"`
	}

	// GoogleCallbackResponse represents the Google OAuth callback response
	GoogleCallbackResponse struct {
		User models.Users `json:"user"`
		Err  string       `json:"err,omitempty"`
	}

	// ServeRequest represents the serve request
	ServeRequest struct {
		Token string `json:"token"`
	}

	// ServeResponse represents the serve response
	ServeResponse struct {
		User models.Users `json:"user"`
		Err  string       `json:"err,omitempty"`
	}
)

// MakeLoginEndpoint creates the login endpoint
func MakeLoginEndpoint(s AuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(LoginRequest)
		redirectURL, err := s.Login(ctx, req.User)
		if err != nil {
			return LoginResponse{Err: err.Error()}, nil
		}
		return LoginResponse{RedirectURL: redirectURL}, nil
	}
}

// MakeRegisterEndpoint creates the register endpoint
func MakeRegisterEndpoint(s AuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(RegisterRequest)
		err := s.Register(ctx, req.User)
		if err != nil {
			return RegisterResponse{Err: err.Error()}, nil
		}
		return RegisterResponse{Message: "User registered successfully"}, nil
	}
}

// MakeGoogleCallbackEndpoint creates the Google callback endpoint
func MakeGoogleCallbackEndpoint(s AuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GoogleCallbackRequest)
		user, err := s.GoogleCallback(ctx, req.Code, req.State)
		if err != nil {
			return GoogleCallbackResponse{Err: err.Error()}, nil
		}
		return GoogleCallbackResponse{User: user}, nil
	}
}

// MakeServeEndpoint creates the serve endpoint
func MakeServeEndpoint(s AuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ServeRequest)
		user, err := s.Serve(ctx, req.Token)
		if err != nil {
			return ServeResponse{Err: err.Error()}, nil
		}
		return ServeResponse{User: user}, nil
	}
}

// NewEndpoints creates all endpoints for the service
func NewEndpoints(s AuthService) Endpoints {
	return Endpoints{
		LoginEndpoint:        MakeLoginEndpoint(s),
		RegisterEndpoint:      MakeRegisterEndpoint(s),
		GoogleCallbackEndpoint: MakeGoogleCallbackEndpoint(s),
		ServeEndpoint:        MakeServeEndpoint(s),
	}
}