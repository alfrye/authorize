package api

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/alfrye/authorize/internal/authorize"
	"github.com/alfrye/authorize/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type (
	AuthHandler interface {
		Login() http.HandlerFunc
		RegisterUsers() http.HandlerFunc
		Serve() http.HandlerFunc
		GoogleReceive() http.HandlerFunc
	}

	Handler struct {
		authService authorize.AuthService
	}
)

func NewAuthHandler(authService authorize.AuthService) AuthHandler {
	return &Handler{
		authService: authService,
	}
}

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("../../client/templates/*.gohtml"))
}

// Login handles the requst to log a user in
func (h *Handler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//TODO: Call the
		if r.Method == http.MethodPost {

			// Get data from post form
			//check for valid password
			// Generate token
			// Create a cookie for to store token
			var user models.Users
			//var t interface{}
			json.NewDecoder(r.Body).Decode(&user)

			//TODO: refactor the login all for each provider
			url, err := h.authService.AuthProvider.Login("")
			if err != nil {
				log.Print("Could not Oauth providers redirect url users")
				return
			}

			// Redirects to OAuth provider login
			http.Redirect(w, r, url, http.StatusSeeOther)

			// Get user from persistence
			// persist data for users
			//?? what db mongo, key/value, cockroach, rdbms

			// retrieveUser, err := h.authService.AuthRepository.GetUser(user.Name)
			// if err != nil {
			// 	log.Print("Could not retrieve users")
			// }

			// result := bcrypt.CompareHashAndPassword([]byte(retrieveUser.Password), []byte(user.Password))

			// if result != nil {
			// 	fmt.Println(result)
			// 	// pass word not valid
			// 	w.WriteHeader(http.StatusUnauthorized)
			// 	//				http.Redirect(w, r, "localhost/authorize/v1/user", http.StatusPermanentRedirect)

			// }

			// call generate token
			//			token := retrieveUser.GenerateToken()

			// Generate Cookie

			// 	cookie := http.Cookie{Name: "auth", Value: token}

			// 	http.SetCookie(w, &cookie)
			// 	fmt.Printf("UserName:%s", result)
			// 	w.WriteHeader(http.StatusOK)
			// 	w.Write([]byte("Login Succeeded"))

		}

	}
}

// GoogleReceive is the oauth callback function
func (h *Handler) GoogleReceive() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		state := r.FormValue("state")

		if state == "" {
			return
		}
		// got code from oauth server
		code := r.FormValue("code")

		if code == "" {
			return
		}

		clt, err := h.authService.AuthProvider.GetOAuthClient(code, r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		resp, err := clt.Get("https://www.googleapis.com/oauth2/v3/userinfo")
		if err != nil {
			log.Printf("Error getting info from Google:%v\n", err.Error())

		}

		defer resp.Body.Close()
		data, _ := ioutil.ReadAll(resp.Body)

		fmt.Printf("Data from Provider:%v\n", string(data))
		user, err := h.authService.AuthProvider.ProcessUserData(data)

		if err != nil {
			log.Println(err)
		}
		// Check to see if the user is already registered
		existingUsers, err := h.authService.AuthRepository.GetUser(user.Name)
		if err != nil {
			log.Println("Error retreiving users from database")
		}
		if (models.Users{}) == existingUsers {
			h.authService.RegisterUser(user)
		}

		// createSession
		h.authService.CreateSession(existingUsers, w)

		//    http.Redirect(w, r, "/", http.StatusSeeOther)

		fmt.Printf("Sending user data to users page:%v", existingUsers)
		tpl.ExecuteTemplate(w, "users.gohtml", existingUsers)

		fmt.Printf("code from provider:%v\n", code)

	}

}

// RegisterUsers handles the requst to register a user
func (h *Handler) RegisterUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodPost {

			// data
			hashedPass, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), bcrypt.MinCost)
			if err != nil {
				fmt.Println(err)
			}

			data := models.Users{
				Name:     r.FormValue("user"),
				Password: string(hashedPass),
				Email:    r.FormValue("email"),
			}
			fmt.Printf("UserName:%v", data)

			// persist data for users
			//?? what db mongo, key/value, cockroach, sqllite rdbms

			h.authService.AuthRepository.CreateUser(data)
			w.WriteHeader(http.StatusCreated)
			_, err = w.Write([]byte("User has been registered"))
			if err != nil {
				fmt.Println(err)
			}
			return
		}

	}
}

// Serve handles the requst to register a user
func (h *Handler) Serve() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check for cookie
		c, err := r.Cookie("auth")
		//	var tpl *template.Template
		//	tpl = template.Must(template.ParseGlob("../../client/templates/index.gohtml"))
		if err != nil {
			log.Println("Cookie does not exist:%v", err)
			c = &http.Cookie{
				Name:  "sessionID",
				Value: "",
			}

		}

		if c.Value == "" {
			pname := h.authService.AuthProvider.GetName()
			tpl.ExecuteTemplate(w, "index.gohtml", pname)
			return
		}

		username, err := h.authService.ParseToken(c.Value)
		if err != nil {
			log.Println("parse token in index route: %v", err)
		}

		currentUser, err := h.authService.AuthRepository.GetUser(username)
		//_, err = h.authService.AuthRepository.GetUser(username)
		if err != nil {
			log.Print("Could not retrieve users")
		}
		//	var tpl *template.Template

		//	tpl = template.Must(template.ParseGlob("../../client/templates/index.gohtml"))

		tpl.ExecuteTemplate(w, "users.gohtml", currentUser)
	}
}
