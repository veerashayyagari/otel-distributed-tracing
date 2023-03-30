package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	m "github.com/veerashayyagari/go-otel/services/models"
)

var userStore = map[string]m.User{
	"1": {ID: "1", Name: "Bill", Email: "bill@gmail.com"},
	"2": {ID: "2", Name: "Jim", Email: "jim@outlook.com"},
	"3": {ID: "3", Name: "Sally", Email: "sally@yahoo.com"},
	"4": {ID: "4", Name: "Mike", Email: "mike@icloud.com"},
	"5": {ID: "5", Name: "Jenni", Email: "jenni@google.com"},
}

func UserAPIRouter() *httprouter.Router {
	router := httprouter.New()
	router.GET("/api/users", getUsersHandler)
	router.GET("/api/users/:id", getUserByIdHandler)
	return router
}

func getUsersHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var users []m.User
	for _, v := range userStore {
		users = append(users, v)
	}

	respond("getUsersHandler", w, http.StatusOK, users)
}

func getUserByIdHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	handler := "getUserByIdHandler"
	if id == "" {
		respond(handler, w, http.StatusBadRequest, fmt.Errorf("not a valid request, could not parse userid"))
		return
	}

	user, ok := userStore[id]
	if !ok {
		respond(handler, w, http.StatusNotFound, fmt.Errorf("user with id(%s) not found", id))
		return
	}

	respond(handler, w, http.StatusOK, user)
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
