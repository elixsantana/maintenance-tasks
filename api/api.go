package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"maintenance-tasks/manager"
	"maintenance-tasks/storage"
)

type Handler struct {
	manager     manager.Manager
	taskPathRgx *regexp.Regexp
}

func CreateHandler(manager manager.Manager) (*Handler, error) {
	taskPathRgx := regexp.MustCompile(`^/task`)
	return &Handler{
		manager:     manager,
		taskPathRgx: taskPathRgx,
	}, nil
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s %s\n", r.Method, r.URL.Path)
	if r.Header.Get("Role") == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	switch r.URL.Path {
	case "/":
		fmt.Fprint(w, "Hello World from an API!")
	case "/tasks":
		h.handleTasks(w, r)
	default:
		if h.taskPathRgx.MatchString(r.URL.Path) {
			h.handleTask(w, r)
		} else {
			http.Error(w, "Not found", http.StatusNotFound)
		}
	}
}

func (h *Handler) handleTasks(w http.ResponseWriter, r *http.Request) {
	isManager := isManager(r.Header.Get("Role"))
	if !isManager {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getTasks(w, isManager)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) handleTask(w http.ResponseWriter, r *http.Request) {
	isManager := isManager(r.Header.Get("Role"))

	switch r.Method {
	case http.MethodGet:
		h.getTask(w, r, isManager)
	case http.MethodPost:
		if isManager {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		h.createTask(w, r)
	case http.MethodPut:
		if isManager {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		h.updateTask(w, r)
	case http.MethodDelete:
		if !isManager {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		h.deleteTask(w, r, isManager)
	default:
		http.Error(w, "Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) getTasks(w http.ResponseWriter, isManager bool) {
	if !isManager {
		fmt.Println("Not a manager")
		http.Error(w, "Failed", http.StatusForbidden)
		return
	}

	tasks, err := h.manager.GetAllTasks()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to retrieve tasks from database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
}

func (h *Handler) getTask(w http.ResponseWriter, r *http.Request, isManager bool) {
	var dTask storage.Task
	var techIdStr string
	var techId int
	var err error

	if !isManager {
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

	task, err := h.manager.GetTask(task_id, techId, isManager)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to retrieve task from database", http.StatusInternalServerError)
		return
	}
	if task == dTask {
		w.WriteHeader(http.StatusOK)
	} else {
		if !isManager {
			h.manager.ExecuteNotification(task, "RETRIEVE TASK")
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(task)
	}
}

func (h *Handler) createTask(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	var task storage.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.validateTask(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.manager.CreateTask(task.Summary, task.TechnicianID, task.Role, now)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to create task", http.StatusInternalServerError)
		return
	}
	task.ID = 1
	task.Performed_date = now
	h.manager.ExecuteNotification(task, "CREATE TASK")
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) updateTask(w http.ResponseWriter, r *http.Request) {
	techIdStr := r.Header.Get("TechId")
	techId, err := strconv.Atoi(techIdStr)
	if err != nil {
		http.Error(w, "Invalid", http.StatusBadRequest)
		return
	}
	var task storage.Task
	err = json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if techId != task.TechnicianID {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	if task.Summary == "" {
		http.Error(w, "Task summary is required", http.StatusBadRequest)
		return
	}

	task, err = h.manager.UpdateTask(task)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.manager.ExecuteNotification(task, "UPDATE TASK")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

func (h *Handler) deleteTask(w http.ResponseWriter, r *http.Request, isManager bool) {
	if !isManager {
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

	err = h.manager.DeleteTask(task_id)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed", http.StatusForbidden)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) validateTask(task *storage.Task) error {
	if task.Summary == "" {
		return errors.New("Summary is required")
	}

	if task.TechnicianID < 1 {
		return errors.New("Technician ID is required")
	}

	if task.Role == "" {
		return errors.New("Role is required")
	}

	return nil
}

func isManager(role string) bool {
	return strings.ToLower(role) == "manager"
}
