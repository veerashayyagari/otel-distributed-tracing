package router

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/veerashayyagari/go-otel/request"
	m "github.com/veerashayyagari/go-otel/services/models"
	"github.com/veerashayyagari/go-otel/tracer"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var templates *template.Template

type SaleDetails struct {
	m.Sale
	ProductName string
}

type UsersPage struct {
	Title string
	Users []m.User
}

type UserSalesPage struct {
	Title string
	Sales []SaleDetails
}

type Router struct {
	http.Handler
	trace.Tracer
}

func init() {
	fs := os.DirFS("./tmpl/")
	templates = template.Must(template.ParseFS(fs, "*.html"))
}

func New(tr trace.Tracer) *Router {
	r := &Router{
		Tracer: tr,
	}
	router := httprouter.New()
	router.GET("/users", tracer.Wrap(r.renderUsersTemplate, tr))
	router.GET("/usersales/:uid", tracer.Wrap(r.renderUserSalesTemplate, tr))
	router.NotFound = http.RedirectHandler("/users", http.StatusMovedPermanently)
	r.Handler = router
	return r
}

func (ro *Router) renderUsersTemplate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	users, err := ro.getUsers(r.Context())
	if err != nil {
		log.Println("getUsers: ", err)
		http.Error(w, "error fetching users", http.StatusInternalServerError)
		return
	}

	data := UsersPage{
		Title: "List Users",
		Users: users,
	}
	err = templates.ExecuteTemplate(w, "users.html", data)
	if err != nil {
		log.Println("executing users.html template with data. ", data, "error", err)
		http.Error(w, "error rendering", http.StatusInternalServerError)
	}
}

func (ro *Router) renderUserSalesTemplate(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	sales, err := ro.getUserSales(r.Context(), p.ByName("uid"))
	if err != nil {
		log.Println("getUserSales: ", "userId", p.ByName("uid"), err)
		http.Error(w, "error fetching usersales", http.StatusInternalServerError)
		return
	}

	salesDetails := make([]SaleDetails, 0, len(sales))

	for _, sale := range sales {
		if p, err := ro.getProductDetails(r.Context(), sale.ProductID); err != nil {
			log.Println("fetching product details", " product: ", sale.ProductID, " error: ", err)
			salesDetails = append(salesDetails, SaleDetails{Sale: sale, ProductName: "Error Fetching Product Name"})
		} else {
			salesDetails = append(salesDetails, SaleDetails{Sale: sale, ProductName: p.Name})
		}
	}

	data := UserSalesPage{
		Title: "List User Sales",
		Sales: salesDetails,
	}

	err = templates.ExecuteTemplate(w, "usersales.html", data)
	if err != nil {
		log.Println("executing usersales.html template with data. ", data, "error", err)
		http.Error(w, "error rendering", http.StatusInternalServerError)
	}
}

func (ro *Router) getUsers(ctx context.Context) ([]m.User, error) {
	_, span := ro.Tracer.Start(ctx, "getUsers")
	startTime := time.Now().UTC()
	defer span.SetAttributes(attribute.Int("execution.time", int(time.Now().UTC().Sub(startTime))))
	defer span.End()

	apiHost, ok := os.LookupEnv("USER_API_URI")
	if !ok {
		return nil, fmt.Errorf("USER_API_URI not found")
	}

	r, err := request.Send(ctx, http.MethodGet, fmt.Sprintf("%s/api/users", apiHost), nil)
	if err != nil {
		return nil, err
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

func (ro *Router) getUserSales(ctx context.Context, uid string) ([]m.Sale, error) {
	_, span := ro.Tracer.Start(ctx, fmt.Sprintf("getUserSales:%s", uid))
	startTime := time.Now().UTC()
	defer span.SetAttributes(attribute.Int("execution.time", int(time.Now().UTC().Sub(startTime))))
	defer span.End()

	apiHost, ok := os.LookupEnv("SALES_API_URI")
	if !ok {
		return nil, fmt.Errorf("SALES_API_URI not found")
	}

	r, err := request.Send(ctx, http.MethodGet, fmt.Sprintf("%s/api/usersales/%s", apiHost, uid), nil)
	if err != nil {
		return nil, err
	}

	var sales []m.Sale
	dec := json.NewDecoder(r.Body)
	err = dec.Decode(&sales)

	if err != nil {
		return nil, fmt.Errorf("decoding response body: %w", err)
	}

	return sales, nil
}

func (ro *Router) getProductDetails(ctx context.Context, id string) (m.Product, error) {
	_, span := ro.Tracer.Start(ctx, fmt.Sprintf("getProductDetails:%s", id))
	startTime := time.Now().UTC()
	defer span.SetAttributes(attribute.Int("execution.time", int(time.Now().UTC().Sub(startTime))))
	defer span.End()

	apiHost, ok := os.LookupEnv("PRODUCT_API_URI")
	if !ok {
		return m.Product{}, fmt.Errorf("PRODUCT_API_URI not found")
	}

	r, err := request.Send(ctx, http.MethodGet, fmt.Sprintf("%s/api/products/%s", apiHost, id), nil)
	if err != nil {
		return m.Product{}, err
	}

	var product m.Product
	dec := json.NewDecoder(r.Body)
	err = dec.Decode(&product)

	if err != nil {
		return m.Product{}, fmt.Errorf("decoding response body: %w", err)
	}

	return product, nil
}
