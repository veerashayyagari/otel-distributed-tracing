package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

type Sale struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	ProductID string    `json:"product_id"`
	Quantity  int       `json:"quantity"`
	SalePrice float64   `json:"sale_price"`
	SaleDate  time.Time `json:"sale_date"`
}

var saleStore = map[string]Sale{
	"100": {"100", "2", "1000", 2, 120.00, time.Now().UTC().AddDate(0, -2, -1)},
	"102": {"102", "2", "1002", 3, 180.30, time.Now().UTC().AddDate(0, -3, -2)},
	"300": {"300", "1", "1001", 2, 150.00, time.Now().UTC().AddDate(0, -1, -2)},
	"204": {"204", "1", "1001", 2, 150.00, time.Now().UTC().AddDate(0, -3, -3)},
	"350": {"350", "3", "1004", 5, 150.00, time.Now().UTC().AddDate(0, -5, -10)},
	"150": {"150", "4", "1003", 10, 1500.00, time.Now().UTC().AddDate(0, -10, -1)},
	"160": {"160", "4", "1005", 5, 100.00, time.Now().UTC().AddDate(0, -5, -1)},
	"250": {"250", "5", "1006", 10, 1600.00, time.Now().UTC().AddDate(0, -10, -5)},
	"360": {"360", "4", "1000", 5, 300.00, time.Now().UTC().AddDate(0, -5, -1)},
	"205": {"205", "3", "1000", 3, 180.00, time.Now().UTC().AddDate(0, -5, -10)},
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

	usrSales := []Sale{}

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
