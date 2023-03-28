package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Product struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

var productStore = map[string]Product{
	"1000": {"1000", "Product A", 60.00},
	"1001": {"1001", "Product B", 75.00},
	"1002": {"1002", "Product C", 60.10},
	"1003": {"1003", "Product D", 150.00},
	"1004": {"1004", "Product E", 30.00},
	"1005": {"1005", "Product F", 20.00},
	"1006": {"1006", "Product G", 160.00},
}

func ProductAPIRouter() *httprouter.Router {
	router := httprouter.New()
	router.GET("/api/products", getAllProducts)
	router.GET("/api/products/:id", getProductsById)
	return router
}

func getAllProducts(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	respond("getAllProducts", w, http.StatusOK, productStore)
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
