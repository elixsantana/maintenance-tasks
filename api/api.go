package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
)

type Handler struct {
}

func CreateHandler() (*Handler, error) {
	return &Handler{}, nil
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	fmt.Println(r.URL)
	// role := r.Header.Get("Role")
	// if role == "" {
	// 	http.Error(w, "Unauthorized", http.StatusUnauthorized)
	// 	return
	// }

	// name := r.Header.Get("Name")
	// if name == "" {
	// 	http.Error(w, "Name not specified", http.StatusInternalServerError)
	// 	return
	// }

	if r.URL.Path == "/" {
		fmt.Fprintf(w, "Hello World from an API!")
		fmt.Println("API hit")
	} else if r.URL.Path == "/tasks" {
		switch r.Method {
		case http.MethodGet:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(`{"Role":"Manager","Task":"Hello","Time":123456789}`)
		case http.MethodPost:
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "Created")
		case http.MethodDelete:
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "Deleted")
		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

	} else if pathRegex := regexp.MustCompile("^/task/(([0-9]+))?$"); pathRegex.MatchString(r.URL.Path) {
		switch r.Method {
		case http.MethodPut:
			pathMatches := pathRegex.FindStringSubmatch(r.URL.Path)
			if len(pathMatches) < 3 || pathMatches[2] == "" {
				http.Error(w, "Invalid task ID", http.StatusBadRequest)
				return
			}
			fmt.Println(pathMatches[2])
			taskID, err := strconv.Atoi(pathMatches[2])
			if err != nil {
				http.Error(w, "Invalid task ID", http.StatusBadRequest)
				return
			}

			w.WriteHeader(http.StatusOK)
			fmt.Println(taskID)
			fmt.Fprintf(w, "Updated")

		case http.MethodDelete:
		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

	} else {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

}
