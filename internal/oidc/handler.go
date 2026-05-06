package oidc

func (p *Provider) setupRoutes() {
	p.router.HandleFunc("/.well-known/openid-configuration", p.handleDiscovery).Methods("GET")
	p.router.HandleFunc("/.well-known/jwks.json", p.handleJWKS).Methods("GET")
	p.router.HandleFunc("/authorize", p.handleAuthorize).Methods("GET", "POST")
	p.router.HandleFunc("/token", p.handleToken).Methods("POST")
	p.router.HandleFunc("/introspect", p.handleIntrospect).Methods("POST")
	p.router.HandleFunc("/userinfo", p.handleUserInfo).Methods("GET")
}
