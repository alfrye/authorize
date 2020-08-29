package api

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/alfrye/authorize/internal/authorize"
	"github.com/alfrye/authorize/internal/models"

	//persistence "github.com/alfrye/authorize/internal/persistence/mongo"
	"golang.org/x/crypto/bcrypt"
)

type (
	AuthHandler interface {
		Login() http.HandlerFunc
		RegisterUsers() http.HandlerFunc
		Serve() http.HandlerFunc
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

			// userName := r.FormValue("user")
			// password := r.FormValue("password")

			// Get user from persistence
			// persist data for users
			//?? what db mongo, key/value, cockroach, sqllite rdbms
			retrieveUser, err := h.authService.AuthRepository.GetUser(user.Name)
			if err != nil {
				log.Print("Could not retrieve users")
			}
			// db := persistence.Database{
			// 	Conn: &persistence.Connection{
			// 		Name: "myconnection",
			// 		Host: "localhost",
			// 		Port: "27017",
			// 	},
			// }

			// 	db.Connect()
			// retrieveUser := db.GetUser(user.Name)

			result := bcrypt.CompareHashAndPassword([]byte(retrieveUser.Password), []byte(user.Password))

			if result != nil {
				fmt.Println(result)
				// pass word not valid
				w.WriteHeader(http.StatusUnauthorized)
				//				http.Redirect(w, r, "localhost/authorize/v1/user", http.StatusPermanentRedirect)

			}

			// call generate token
			token := retrieveUser.GenerateToken()

			// Generate Cookie

			cookie := http.Cookie{Name: "auth", Value: token}

			http.SetCookie(w, &cookie)
			fmt.Printf("UserName:%s", result)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Login Succeeded"))

		}

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
			fmt.Printf("UserName:%s", data)

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

		var tpl *template.Template

		tpl = template.Must(template.ParseGlob("../../client/templates/*"))

		tpl.ExecuteTemplate(w, "index.gohtml", nil)
	}
}
