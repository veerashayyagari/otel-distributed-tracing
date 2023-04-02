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

var userStore = map[string]m.User{
	"1": {ID: "1", Name: "Bill", Email: "bill@gmail.com"},
	"2": {ID: "2", Name: "Jim", Email: "jim@outlook.com"},
	"3": {ID: "3", Name: "Sally", Email: "sally@yahoo.com"},
	"4": {ID: "4", Name: "Mike", Email: "mike@icloud.com"},
	"5": {ID: "5", Name: "Jenni", Email: "jenni@google.com"},
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
	router.GET("/api/users", tracer.Wrap(getUsersHandler, tr))
	router.GET("/api/users/:id", tracer.Wrap(getUserByIDHandler, tr))
	r.Handler = router
	return r
}

func getUsersHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var users []m.User
	for _, v := range userStore {
		users = append(users, v)
	}

	response.Send("getUsersHandler", w, http.StatusOK, users)
}

func getUserByIDHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	handler := "getUserByIdHandler"
	if id == "" {
		response.Send(handler, w, http.StatusBadRequest, fmt.Errorf("not a valid request, could not parse userid"))
		return
	}

	user, ok := userStore[id]
	if !ok {
		response.Send(handler, w, http.StatusNotFound, fmt.Errorf("user with id(%s) not found", id))
		return
	}

	response.Send(handler, w, http.StatusOK, user)
}
