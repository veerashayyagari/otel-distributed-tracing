package response

import (
	"encoding/json"
	"net/http"
)

func Send(handler string, w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	buf, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error fetching user(s)"))
	} else {
		w.WriteHeader(statusCode)
		w.Write(buf)
	}
}
