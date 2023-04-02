package router

import (
	"fmt"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/veerashayyagari/go-otel/response"
	m "github.com/veerashayyagari/go-otel/services/models"
	"github.com/veerashayyagari/go-otel/tracer"
	"go.opentelemetry.io/otel/trace"
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

type Router struct {
	http.Handler
	trace.Tracer
}

func New(tr trace.Tracer) *Router {
	r := &Router{
		Tracer: tr,
	}
	router := httprouter.New()
	router.GET("/api/sales/:id", tracer.Wrap(getSalesByID, tr))
	router.GET("/api/usersales/:uid", tracer.Wrap(getSalesByUserID, tr))
	r.Handler = router
	return r
}

func getSalesByID(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	handler := "getSalesById"
	if id == "" {
		response.Send(handler, w, http.StatusBadRequest, fmt.Errorf("not a valid request, could not parse saleid"))
		return
	}

	sale, ok := saleStore[id]
	if !ok {
		response.Send(handler, w, http.StatusNotFound, fmt.Errorf("sale with id(%s) not found", id))
		return
	}

	response.Send(handler, w, http.StatusOK, sale)
}

func getSalesByUserID(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("uid")
	handler := "getSalesByUserId"
	if id == "" {
		response.Send(handler, w, http.StatusBadRequest, fmt.Errorf("not a valid request, could not parse userid"))
		return
	}

	usrSales := []m.Sale{}

	for _, sale := range saleStore {
		if sale.UserID == id {
			usrSales = append(usrSales, sale)
		}
	}

	response.Send(handler, w, http.StatusOK, usrSales)
}
