package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"time"

	"github.com/julienschmidt/httprouter"
	m "github.com/veerashayyagari/go-otel/services/models"
	"github.com/veerashayyagari/go-otel/tracer"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
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

type App struct {
	http.Handler
	trace.Tracer
}

func init() {
	fs := os.DirFS("./tmpl/")
	templates = template.Must(template.ParseFS(fs, "*.html"))
}

func New(tr trace.Tracer) *App {
	a := &App{
		Tracer: tr,
	}
	router := httprouter.New()
	router.NotFound = http.RedirectHandler("/users", http.StatusMovedPermanently)
	router.GET("/users", tracer.Wrap(a.renderUsersTemplate, tr))
	router.GET("/users/:uid", tracer.Wrap(a.renderUserSalesTemplate, tr))
	a.Handler = router
	return a
}

func (a *App) renderUsersTemplate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	users, err := a.getUsers(r.Context())
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

func (a *App) renderUserSalesTemplate(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	sales, err := a.getUserSales(r.Context(), p.ByName("uid"))
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

func (a *App) getUsers(ctx context.Context) ([]m.User, error) {
	_, span := a.Tracer.Start(ctx, "getUsers")
	startTime := time.Now().UTC()
	defer span.SetAttributes(attribute.Int("execution.time", int(time.Now().UTC().Sub(startTime))))
	defer span.End()

	apiHost, ok := os.LookupEnv("USER_API_URI")
	if !ok {
		return nil, fmt.Errorf("USER_API_URI not found")
	}

	r, err := http.Get(fmt.Sprintf("%s/api/users", apiHost))
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

func (a *App) getUserSales(ctx context.Context, uid string) ([]m.Sale, error) {
	_, span := a.Tracer.Start(ctx, fmt.Sprintf("getUserSales:%s", uid))
	startTime := time.Now().UTC()
	defer span.SetAttributes(attribute.Int("execution.time", int(time.Now().UTC().Sub(startTime))))
	defer span.End()

	apiHost, ok := os.LookupEnv("SALES_API_URI")
	if !ok {
		return nil, fmt.Errorf("SALES_API_URI not found")
	}

	r, err := http.Get(fmt.Sprintf("%s/api/usersales/%s", apiHost, uid))
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
