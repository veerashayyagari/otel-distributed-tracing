package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	m "github.com/veerashayyagari/go-otel/services/models"
)

var saleStore = map[string]m.Sale{
	"100": {ID: "100", UserID: "2", ProductID: "1000", Quantity: 2, SalePrice: 120.00, SaleDate: time.Now().UTC().AddDate(0, -2, -1)},
	"102": {ID: "102", UserID: "2", ProductID: "1002", Quantity: 3, SalePrice: 180.30, SaleDate: time.Now().UTC().AddDate(0, -3, -2)},
	"300": {ID: "300", UserID: "1", ProductID: "1001", Quantity: 2, SalePrice: 150.00, SaleDate: time.Now().UTC().AddDate(0, -1, -2)},
	"204": {ID: "204", UserID: "1", ProductID: "1001", Quantity: 2, SalePrice: 150.00, SaleDate: time.Now().UTC().AddDate(0, -3, -3)},
	"350": {ID: "350", UserID: "3", ProductID: "1004", Quantity: 5, SalePrice: 150.00, SaleDate: time.Now().UTC().AddDate(0, -5, -10)},
	"150": {ID: "150", UserID: "4", ProductID: "1003", Quantity: 10, SalePrice: 1500.00, SaleDate: time.Now().UTC().AddDate(0, -10, -1)},
	"160": {ID: "160", UserID: "4", ProductID: "1005", Quantity: 5, SalePrice: 100.00, SaleDate: time.Now().UTC().AddDate(0, -5, -1)},
	"250": {ID: "250", UserID: "5", ProductID: "1006", Quantity: 10, SalePrice: 1600.00, SaleDate: time.Now().UTC().AddDate(0, -10, -5)},
	"360": {ID: "360", UserID: "4", ProductID: "1000", Quantity: 5, SalePrice: 300.00, SaleDate: time.Now().UTC().AddDate(0, -5, -1)},
	"205": {ID: "205", UserID: "3", ProductID: "1000", Quantity: 3, SalePrice: 180.00, SaleDate: time.Now().UTC().AddDate(0, -5, -10)},
}

func SalesAPIRouter() *httprouter.Router {
	router := httprouter.New()
	router.GET("/api/sales/:id", getSalesById)
	router.GET("/api/usersales/:uid", getSalesByUserId)
	return router
}

func getSalesById(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	handler := "getSalesById"
	if id == "" {
		respond(handler, w, http.StatusBadRequest, fmt.Errorf("not a valid request, could not parse saleid"))
		return
	}

	sale, ok := saleStore[id]
	if !ok {
		respond(handler, w, http.StatusNotFound, fmt.Errorf("sale with id(%s) not found", id))
		return
	}

	respond(handler, w, http.StatusOK, sale)
}

func getSalesByUserId(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("uid")
	handler := "getSalesByUserId"
	if id == "" {
		respond(handler, w, http.StatusBadRequest, fmt.Errorf("not a valid request, could not parse userid"))
		return
	}

	usrSales := []m.Sale{}

	for _, sale := range saleStore {
		if sale.UserID == id {
			usrSales = append(usrSales, sale)
		}
	}

	respond(handler, w, http.StatusOK, usrSales)
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
