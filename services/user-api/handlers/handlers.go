package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var userStore = map[string]User{
	"1": {"1", "Bill", "bill@gmail.com"},
	"2": {"2", "Jim", "jim@outlook.com"},
	"3": {"3", "Sally", "sally@yahoo.com"},
	"4": {"4", "Mike", "mike@icloud.com"},
	"5": {"5", "Jenni", "jenni@google.com"},
}

func UserAPIRouter() *httprouter.Router {
	router := httprouter.New()
	router.GET("/api/users", getUsersHandler)
	router.GET("/api/users/:id", getUserByIdHandler)
	return router
}

func getUsersHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	respond("getUsersHandler", w, http.StatusOK, userStore)
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
