package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/alfrye/authorize/internal/authorize"
	"github.com/alfrye/authorize/internal/models"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type (
	googleProvider struct {
		providerName string
	}

	userData struct {
		Name    string
		Email   string
		Picture string
		Sub     string
	}
)

var conf = &oauth2.Config{
	ClientID:     "290625463247-dcqgtr6t69tkc5qet9rbubjj1gi8nc00.apps.googleusercontent.com",
	ClientSecret: "ihfjio6HajQ1f3JbQcXgA8ev",
	Endpoint:     google.Endpoint,
	RedirectURL:  "http://localhost:9010/authorize/v1/oauth/google",
	Scopes: []string{
		"profile",
		"email",
		//"https://www.googleapis.com/auth/userinfo.email",
	},
}

// This will provide the implementation of an authorization provider

// NewLocalAuthProvider creates a instance of a local auth provider
func NewGoogleAuthProvider() (authorize.AuthProvider, error) {
	ap := &googleProvider{
		providerName: "google",
	}

	return ap, nil

}

func (p *googleProvider) GetName() string {
	return p.providerName
}

// Login implementation of the login method
func (p *googleProvider) Login(username string) (string, error) {

	state := uuid.NewV4()

	url := conf.AuthCodeURL(state.String())
	//	http.Redirect(w, r, url, http.StatusSeeOther)

	fmt.Println(state)

	return url, nil
}

// RegisterUsers is an implementation of the Register users
func (p *googleProvider) Register(users models.Users) error {
	return nil
}

func (p *googleProvider) GetOAuthClient(code string, ctx context.Context) (*http.Client, error) {
	//Call auth provider for the rest of this
	//get token from code
	t, err := conf.Exchange(ctx, code)
	if err != nil {
		return nil, err

	}

	ts := conf.TokenSource(ctx, t)

	// Create new client
	clt := oauth2.NewClient(ctx, ts)

	return clt, nil
}

func (p *googleProvider) ProcessUserData(data []byte) (models.Users, error) {
	users := models.Users{}
	u := userData{}
	jerr := json.Unmarshal(data, &u)
	if jerr != nil {
		log.Printf("Error unmarshaling data info from Google:%v\n", jerr.Error())
		return models.Users{}, jerr

	}
	fmt.Printf("Data for user%v\n", u)
	users.Name = u.Name
	users.AvatarURL = u.Picture
	users.Email = u.Email
	users.UserID = u.Sub
	return users, nil

}
