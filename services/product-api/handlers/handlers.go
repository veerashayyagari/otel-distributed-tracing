package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	m "github.com/veerashayyagari/go-otel/services/models"
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

func ProductAPIRouter() *httprouter.Router {
	router := httprouter.New()
	router.GET("/api/products", getAllProducts)
	router.GET("/api/products/:id", getProductsById)
	return router
}

func getAllProducts(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var products []m.Product
	for _, v := range productStore {
		products = append(products, v)
	}
	respond("getAllProducts", w, http.StatusOK, products)
}

func getProductsById(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	handler := "getProductsById"
	if id == "" {
		respond(handler, w, http.StatusBadRequest, fmt.Errorf("not a valid request, could not parse productid"))
		return
	}

	product, ok := productStore[id]
	if !ok {
		respond(handler, w, http.StatusNotFound, fmt.Errorf("product with id(%s) not found", id))
		return
	}

	respond(handler, w, http.StatusOK, product)
}

func respond(handler string, w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	buf, err := json.Marshal(data)
	if err != nil {
		log.Printf("ERROR: %s marshaling json %s \n", handler, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error fetching user(s)"))
	} else {
		w.WriteHeader(statusCode)
		w.Write(buf)
	}
}
