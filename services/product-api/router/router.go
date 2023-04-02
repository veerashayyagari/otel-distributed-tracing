package router

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/veerashayyagari/go-otel/response"
	m "github.com/veerashayyagari/go-otel/services/models"
	"github.com/veerashayyagari/go-otel/tracer"
	"go.opentelemetry.io/otel/trace"
)

var productStore = map[string]m.Product{
	"1000": {ID: "1000", Name: "Product A", Price: 60.00},
	"1001": {ID: "1001", Name: "Product B", Price: 75.00},
	"1002": {ID: "1002", Name: "Product C", Price: 60.10},
	"1003": {ID: "1003", Name: "Product D", Price: 150.00},
	"1004": {ID: "1004", Name: "Product E", Price: 30.00},
	"1005": {ID: "1005", Name: "Product F", Price: 20.00},
	"1006": {ID: "1006", Name: "Product G", Price: 160.00},
}

type Router struct {
	http.Handler
	trace.Tracer
}

func New(tr trace.Tracer) *Router {
	r := &Router{
		Tracer: tr,
	}
	router := httprouter.New()
	router.GET("/api/products", tracer.Wrap(getAllProducts, tr))
	router.GET("/api/products/:id", tracer.Wrap(getProductsById, tr))
	return r
}

func getAllProducts(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var products []m.Product
	for _, v := range productStore {
		products = append(products, v)
	}
	response.Send("getAllProducts", w, http.StatusOK, products)
}

func getProductsById(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	handler := "getProductsById"
	if id == "" {
		response.Send(handler, w, http.StatusBadRequest, fmt.Errorf("not a valid request, could not parse productid"))
		return
	}

	product, ok := productStore[id]
	if !ok {
		response.Send(handler, w, http.StatusNotFound, fmt.Errorf("product with id(%s) not found", id))
		return
	}

	response.Send(handler, w, http.StatusOK, product)
}
