package transport

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/alfrye/authorize/internal/endpoints"
)

func NewHTTPHandler(endpoints endpoints.Endpoints) http.Handler {
	m := http.NewServeMux()

	m.Handle("/auth/login", kithttp.NewServer(
		endpoints.LoginEndpoint,
		decodeLoginRequest,
		encodeLoginResponse,
	))

	m.Handle("/auth/register", kithttp.NewServer(
		endpoints.RegisterEndpoint,
		decodeRegisterRequest,
		encodeRegisterResponse,
	))

	m.Handle("/auth/google/callback", kithttp.NewServer(
		endpoints.GoogleCallbackEndpoint,
		decodeGoogleCallbackRequest,
		encodeGoogleCallbackResponse,
	))

	m.Handle("/auth/serve", kithttp.NewServer(
		endpoints.ServeEndpoint,
		decodeServeRequest,
		encodeServeResponse,
	))

	return m
}

func decodeLoginRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request endpoints.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&request.User); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeLoginResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func decodeRegisterRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request endpoints.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&request.User); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeRegisterResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func decodeGoogleCallbackRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request endpoints.GoogleCallbackRequest
	
	// Try to get from form data first (for OAuth callbacks)
	if r.Method == http.MethodGet {
		request.Code = r.URL.Query().Get("code")
		request.State = r.URL.Query().Get("state")
	} else {
		// Try JSON body
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			return nil, err
		}
	}
	
	return request, nil
}

func encodeGoogleCallbackResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func decodeServeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request endpoints.ServeRequest
	
	// Try to get token from cookie first
	cookie, err := r.Cookie("auth")
	if err == nil && cookie.Value != "" {
		request.Token = cookie.Value
	} else {
		// Try to get from Authorization header
		authHeader := r.Header.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			request.Token = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			// Try JSON body
			if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
				return nil, err
			}
		}
	}
	
	return request, nil
}

func encodeServeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}