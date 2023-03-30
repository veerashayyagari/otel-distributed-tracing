package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	m "github.com/veerashayyagari/go-otel/services/models"
)

var templates *template.Template

type UsersPage struct {
	Title string
	Users []m.User
}

type UserSalesPage struct {
	Title string
	Sales []m.Sale
}

func init() {
	fs := os.DirFS("./tmpl/")
	templates = template.Must(template.ParseFS(fs, "*.html"))
}

func AppRouter() *httprouter.Router {
	router := httprouter.New()
	router.GET("/users", renderUsersTemplate)
	router.GET("/users/:uid", renderUserSalesTemplate)
	router.NotFound = http.RedirectHandler("/users", http.StatusMovedPermanently)

	return router
}

func renderUsersTemplate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	users, err := getUsers()
	if err != nil {
		fmt.Println("getUsers: ", err)
		http.Error(w, "error fetching users", http.StatusInternalServerError)
		return
	}

	data := UsersPage{
		Title: "List Users",
		Users: users,
	}
	err = templates.ExecuteTemplate(w, "users.html", data)
	if err != nil {
		fmt.Println("executing users.html template with data. ", data, "error", err)
		http.Error(w, "error rendering", http.StatusInternalServerError)
	}
}

func renderUserSalesTemplate(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	sales, err := getUserSales(p.ByName("uid"))
	if err != nil {
		fmt.Println("getUserSales: ", "userId", p.ByName("uid"), err)
		http.Error(w, "error fetching usersales", http.StatusInternalServerError)
		return
	}

	data := UserSalesPage{
		Title: "List User Sales",
		Sales: sales,
	}

	err = templates.ExecuteTemplate(w, "usersales.html", data)
	if err != nil {
		fmt.Println("executing usersales.html template with data. ", data, "error", err)
		http.Error(w, "error rendering", http.StatusInternalServerError)
	}
}

func getUsers() ([]m.User, error) {
	r, err := http.Get("http://localhost:4000/api/users")
	if err != nil {
		return nil, fmt.Errorf("fetching users: %w", err)
	}

	if r.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received status code: %d", r.StatusCode)
	}

	var users []m.User
	dec := json.NewDecoder(r.Body)
	err = dec.Decode(&users)

	if err != nil {
		return nil, fmt.Errorf("decoding response body: %w", err)
	}

	return users, nil
}

func getUserSales(uid string) ([]m.Sale, error) {
	r, err := http.Get(fmt.Sprintf("http://localhost:5000/api/usersales/%s", uid))
	if err != nil {
		return nil, fmt.Errorf("fetching user sales: %w", err)
	}

	if r.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received status code: %d", r.StatusCode)
	}

	var sales []m.Sale
	dec := json.NewDecoder(r.Body)
	err = dec.Decode(&sales)

	if err != nil {
		return nil, fmt.Errorf("decoding response body: %w", err)
	}

	return sales, nil
}
