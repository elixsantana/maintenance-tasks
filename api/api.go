package api

import (
	"encoding/json"
	"fmt"
	"maintenance-tasks/manager"
	database "maintenance-tasks/storage"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type Handler struct {
	m *manager.Manager
}

func CreateHandler(manager *manager.Manager) (*Handler, error) {
	return &Handler{
		m: manager,
	}, nil
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL)
	role := r.Header.Get("Role")
	if role == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if r.URL.Path == "/" {
		fmt.Fprintf(w, "Hello World from an API!")
		fmt.Println("API hit")
	} else if r.URL.Path == "/tasks" {
		switch r.Method {
		case http.MethodGet:
			if !isManager(role) {
				fmt.Println("Not a manager")
				http.Error(w, "Failed", http.StatusForbidden)
				return
			}
			tasks, err := h.m.GetAllTasks()
			if err != nil {
				fmt.Println(err)
				http.Error(w, "Failed to retrieve task from database", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(tasks)
		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
	} else if pathRegex := regexp.MustCompile("^/task"); pathRegex.MatchString(r.URL.Path) {
		switch r.Method {
		case http.MethodGet:
			var techIdStr string
			var techId int
			var err error
			if !isManager(role) {
				techIdStr = r.Header.Get("TechId")
				techId, err = strconv.Atoi(techIdStr)
				if err != nil {
					http.Error(w, "Invalid", http.StatusBadRequest)
					return
				}
			}

			taskId := r.URL.Query().Get("id")
			task_id, err := strconv.Atoi(taskId)
			if err != nil {
				http.Error(w, "Invalid id", http.StatusBadRequest)
				return
			}

			task, err := h.m.GetTask(task_id, techId, isManager(role))
			if err != nil {
				fmt.Println(err)
				http.Error(w, "Failed to retrieve task from database", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(task)
		case http.MethodPost:
			if isManager(role) {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			err := r.ParseForm()
			if err != nil {
				fmt.Println(err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			summary := r.PostFormValue("summary")
			techId := r.PostFormValue("techId")
			role := r.PostFormValue("role")

			if summary == "" || techId == "" || role == "" {
				http.Error(w, "Missing required parameters", http.StatusBadRequest)
				return
			}
			err = h.m.CreateTask(summary, techId, role)
			if err != nil {
				fmt.Println(err)
				http.Error(w, "Failed creating task", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
		// case http.MethodDelete:
		// 	w.WriteHeader(http.StatusOK)
		// 	fmt.Fprintf(w, "Deleted")
		case http.MethodPut:
			if isManager(role) {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			var task database.Task
			err := json.NewDecoder(r.Body).Decode(&task)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if task.Summary == "" {
				http.Error(w, "Task summary is required", http.StatusBadRequest)
				return
			}

			task, err = h.m.UpdateTask(task)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(task)

		case http.MethodDelete:
			if !isManager(role) {
				fmt.Println("Not a manager")
				http.Error(w, "Failed", http.StatusForbidden)
				return
			}

			taskId := r.URL.Query().Get("id")
			task_id, err := strconv.Atoi(taskId)
			if err != nil {
				http.Error(w, "Invalid id", http.StatusBadRequest)
				return
			}

			err = h.m.DeleteTask(task_id)
			if err != nil {
				fmt.Println(err)
				http.Error(w, "Failed", http.StatusForbidden)
				return
			}
			w.WriteHeader(http.StatusOK)

		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

	} else {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

}

func isManager(role string) bool {
	return strings.ToLower(role) == "manager"
}
