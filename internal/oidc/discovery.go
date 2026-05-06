package oidc

import (
	"encoding/json"
	"net/http"
)

func (p *Provider) handleDiscovery(w http.ResponseWriter, r *http.Request) {
	discovery := map[string]interface{}{
		"issuer":                                p.config.Issuer,
		"authorization_endpoint":                p.config.Issuer + "/authorize",
		"token_endpoint":                        p.config.Issuer + "/token",
		"userinfo_endpoint":                     p.config.Issuer + "/userinfo",
		"jwks_uri":                              p.config.Issuer + "/.well-known/jwks.json",
		"introspection_endpoint":                p.config.Issuer + "/introspect",
		"scopes_supported":                      []string{"openid", "profile", "email"},
		"response_types_supported":              []string{"code"},
		"grant_types_supported":                 []string{"authorization_code", "refresh_token"},
		"subject_types_supported":               []string{"public"},
		"id_token_signing_alg_values_supported": []string{"ES256"},
		"code_challenge_methods_supported":      []string{"S256"},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(discovery)
}
