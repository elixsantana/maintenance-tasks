package api

import (
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
	}
}
