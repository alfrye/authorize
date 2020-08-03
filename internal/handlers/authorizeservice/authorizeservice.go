package authorizeservice

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/alfrye/authorize/internal/models"
	"github.com/alfrye/authorize/internal/persistence"
	"golang.org/x/crypto/bcrypt"
)

// Login handles the requst to log a user in
func Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

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
			db := persistence.Database{
				Conn: &persistence.Connection{
					Name: "myconnection",
					Host: "localhost",
					Port: "27017",
				},
			}

			db.Connect()
			retrieveUser := db.GetUser(user.Name)

			result := bcrypt.CompareHashAndPassword([]byte(retrieveUser.Password), []byte(user.Password))

			if result != nil {
				fmt.Println(result)
				// pass word not valid
				w.WriteHeader(http.StatusUnauthorized)
				//	http.Redirect(w, r, "localhost/authorize/v1/user", http.StatusPermanentRedirect)

			}

			// call generate token
			token := retrieveUser.GenereateToken()

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
func RegisterUsers() http.HandlerFunc {
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
			db := persistence.Database{
				Conn: &persistence.Connection{
					Name: "myconnection",
					Host: "localhost",
					Port: "27017",
				},
			}

			dbClient := db.Connect()
			dbClient.CreateUser(data)
			_, err = w.Write([]byte("User has been registered"))
			if err != nil {
				fmt.Println(err)
			}
			return
		}

	}

}

// RegisterUsers handles the requst to register a user
func Serve() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var tpl *template.Template

		tpl = template.Must(template.ParseGlob("../../client/templates/*"))

		tpl.ExecuteTemplate(w, "index.gohtml", nil)
	}
}
