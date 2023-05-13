package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Handler struct {
}

func CreateHandler() (*Handler, error) {
	return &Handler{}, nil
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL)
	if r.URL.Path == "/" {
		fmt.Fprintf(w, "Hello World from an API!")
		fmt.Println("API hit")
	} else if r.URL.Path == "/tasks" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(`{"Role":"Manager","Task":"Hello","Time":123456789}`)
	}

}
