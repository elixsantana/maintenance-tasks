package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"maintenance-tasks/manager"
	"maintenance-tasks/storage"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"gotest.tools/assert"
)

func createRequest(
	mockManager manager.Manager,
	endpoint string,
	verb string,
	rawAuthHeader []string,
	body string,
	t *testing.T,
) *httptest.ResponseRecorder {
	h, _ := CreateHandler(mockManager)

	req, err := http.NewRequest(verb, endpoint, bytes.NewBuffer([]byte(body)))
	assert.NilError(t, err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Role", rawAuthHeader[0])
	req.Header.Set("TechId", rawAuthHeader[1])

	mockHandler := http.HandlerFunc(h.ServeHTTP)

	r := httptest.NewRecorder()
	mockHandler.ServeHTTP(r, req)

	return r
}

func TestApi_Authorized(t *testing.T) {
	headers := []string{"manager", ""}
	r := createRequest(TestManager{}, "/tasks", "GET", headers, "", t)
	assert.Equal(t, r.Code, http.StatusOK)

	headers = []string{"technician", "1"}
	r = createRequest(TestManager{}, "/task?id=3", "GET", headers, "", t)
	assert.Equal(t, r.Code, http.StatusOK)
}
func TestApi_Unauthorized(t *testing.T) {
	headers := []string{"", ""}
	r := createRequest(TestManager{}, "/tasks", "GET", headers, "", t)
	assert.Equal(t, r.Code, http.StatusUnauthorized)
}

func TestApi_RouteNotFound(t *testing.T) {
	headers := []string{"manager", "20"}
	r := createRequest(TestManager{}, "/not/found", "GET", headers, "", t)
	assert.Equal(t, r.Code, http.StatusNotFound)
}

func TestApi_Forbidden(t *testing.T) {
	headers := []string{"technician", ""}
	r := createRequest(TestManager{}, "/tasks", "GET", headers, "", t)
	assert.Equal(t, r.Code, http.StatusForbidden)

	techId := 5
	headers = []string{"technician", "1"}
	body := fmt.Sprintf(`{"id":%d,"summary":"Updated Task","performed_date":"2023-05-14T00:00:00Z","technician_id":999,"role":"technician"}`, techId)
	r = createRequest(TestManager{}, "/task?id=2", "PUT", headers, body, t)
	assert.Equal(t, r.Code, http.StatusForbidden)

	headers = []string{"technician", ""}
	r = createRequest(TestManager{}, "/task?id=1", "DELETE", headers, "", t)
	assert.Equal(t, r.Code, http.StatusForbidden)

	headers = []string{"manager", ""}
	r = createRequest(TestManager{}, "/task", "POST", headers, "", t)
	assert.Equal(t, r.Code, http.StatusForbidden)

	headers = []string{"manager", ""}
	r = createRequest(TestManager{}, "/task", "PUT", headers, "", t)
	assert.Equal(t, r.Code, http.StatusForbidden)
}

func TestApi_BadRequest(t *testing.T) {
	headers := []string{"manager", "20"}
	r := createRequest(TestManager{}, "/task?id=xx", "DELETE", headers, "", t)
	assert.Equal(t, r.Code, http.StatusBadRequest)

	headers = []string{"manager", "20"}
	r = createRequest(TestManager{}, "/task?id=xx", "GET", headers, "", t)
	assert.Equal(t, r.Code, http.StatusBadRequest)

	headers = []string{"technician", "xx"}
	r = createRequest(TestManager{}, "/task?id=1", "GET", headers, "", t)
	assert.Equal(t, r.Code, http.StatusBadRequest)

	task := `{"summaaaaary":"Create Task","tsechnician_id":999,"role":"technician"}`
	r = createRequest(TestManager{}, "/task", "POST", headers, task, t)
	assert.Equal(t, r.Code, http.StatusBadRequest)

	task = `{"summaaaaary":"","tsechnician_id":999,"role":"technician"}`
	r = createRequest(TestManager{}, "/task", "POST", headers, task, t)
	assert.Equal(t, r.Code, http.StatusBadRequest)

	task = `{"summary":"","technician_id":,"role":"technician"}`
	r = createRequest(TestManager{}, "/task", "POST", headers, task, t)
	assert.Equal(t, r.Code, http.StatusBadRequest)

	task = `{"summary":"","technician_id":,"role":"technician"}`
	r = createRequest(TestManager{}, "/task", "PUT", headers, task, t)
	assert.Equal(t, r.Code, http.StatusBadRequest)

	task = `{"summary":"Create Task","technician_id":,"role":"technician"}`
	r = createRequest(TestManager{}, "/task", "POST", headers, task, t)
	assert.Equal(t, r.Code, http.StatusBadRequest)

	task = `{"summary":"Create Task","technician_id":")&@&^#(&^@&@&@),"role":"technician"}`
	r = createRequest(TestManager{}, "/task", "PUT", headers, task, t)
	assert.Equal(t, r.Code, http.StatusBadRequest)

	task = `{"summary":"","technician_id":1,"role":"technician"}`
	r = createRequest(TestManager{}, "/task", "PUT", headers, task, t)
	assert.Equal(t, r.Code, http.StatusBadRequest)

	headers = []string{"technician", "1"}
	task = `{"summarewewy":"Create","technician_id":1,"roewle":"te869T9Hian"}`
	r = createRequest(TestManager{}, "/task", "PUT", headers, task, t)
	assert.Equal(t, r.Code, http.StatusBadRequest)
}

func TestApi_MethodNotAllowed(t *testing.T) {
	headers := []string{"Manager", "20"}
	r := createRequest(TestManager{}, "/tasks", "PUT", headers, "", t)
	assert.Equal(t, r.Code, http.StatusMethodNotAllowed)

	headers = []string{"Manager", "20"}
	r = createRequest(TestManager{}, "/tasks", "POST", headers, "", t)
	assert.Equal(t, r.Code, http.StatusMethodNotAllowed)

	headers = []string{"Manager", "20"}
	r = createRequest(TestManager{}, "/tasks", "PATCH", headers, "", t)
	assert.Equal(t, r.Code, http.StatusMethodNotAllowed)

	headers = []string{"Manager", "20"}
	r = createRequest(TestManager{}, "/task", "PATCH", headers, "", t)
	assert.Equal(t, r.Code, http.StatusMethodNotAllowed)
}

func TestHandler_ServeHTTP_Tasks_Get(t *testing.T) {
	headers := []string{"manager", "20"}

	r := createRequest(TestManager{}, "/tasks", "GET", headers, "", t)
	assert.Equal(t, r.Code, http.StatusOK)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	var actualBody []storage.Task
	_ = json.Unmarshal(body, &actualBody)

	expectedBody := GetMockDataSlice(0)
	assert.DeepEqual(t, actualBody, expectedBody)
}

func TestHandler_ServeHTTP_Task_Get(t *testing.T) {
	taskId := 2
	pathAndQueryParams := fmt.Sprintf("/task?id=%d", taskId)
	headers := []string{"manager", "20"}

	r := createRequest(TestManager{}, pathAndQueryParams, "GET", headers, "", t)
	assert.Equal(t, r.Code, http.StatusOK)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}

	var act storage.Task
	_ = json.Unmarshal(body, &act)
	var actualBody []storage.Task
	actualBody = append(actualBody, act)

	expectedBody := GetMockDataSlice(taskId)
	assert.Equal(t, len(actualBody), len(expectedBody))
	assert.DeepEqual(t, actualBody, expectedBody)
}

func TestHandler_ServeHTTP_Task_Post(t *testing.T) {
	pathAndQueryParams := "/task"
	headers := []string{"technician", "999"}

	task := `{"summary":"Create Task","technician_id":999,"role":"technician"}`

	r := createRequest(TestManager{}, pathAndQueryParams, "POST", headers, task, t)
	assert.Equal(t, r.Code, http.StatusCreated)
}

func TestHandler_ServeHTTP_Task_Put(t *testing.T) {
	taskId := 6
	pathAndQueryParams := "/task"
	headers := []string{"technician", "999"}
	task := fmt.Sprintf(`{"id":%d,"summary":"Updated Task","performed_date":"2023-05-14T00:00:00Z","technician_id":999,"role":"technician"}`, taskId)

	r := createRequest(TestManager{}, pathAndQueryParams, "PUT", headers, task, t)
	assert.Equal(t, r.Code, http.StatusOK)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}

	actualBody := strings.TrimSpace(string(body))
	expectedBody := strings.TrimSpace(task)

	assert.DeepEqual(t, actualBody, expectedBody)
}

func TestHandler_ServeHTTP_Task_Delete(t *testing.T) {
	id := 6
	pathAndQueryParams := fmt.Sprintf("/task?id=%d", id)
	headers := []string{"Manager", "20"}
	tManager := TestManager{}
	originalMockData := mockData

	r := createRequest(tManager, pathAndQueryParams, "DELETE", headers, "", t)

	actual, _ := tManager.GetAllTasks()

	delete(originalMockData, id)
	expected := OrderedTasksSlice(originalMockData)

	assert.Equal(t, r.Code, http.StatusOK)
	assert.DeepEqual(t, actual, expected)
}

type TestManager struct{}

func (t TestManager) Start() {

}
func (t TestManager) Stop() {

}
func (t TestManager) GetAllTasks() ([]storage.Task, error) {
	return GetMockDataSlice(0), nil
}
func (t TestManager) CreateTask(summary string, techId int, role string, now time.Time) error {
	return nil
}
func (t TestManager) UpdateTask(task storage.Task) (storage.Task, error) {
	return ModifyTaskMockData(task), nil
}
func (t TestManager) GetTask(task_id int, tech_id int, manager bool) (storage.Task, error) {
	task := GetMockDataSlice(task_id)
	return task[0], nil
}
func (t TestManager) DeleteTask(taskID int) error {
	delete(mockData, taskID)
	return nil
}
func (t TestManager) ReceiveNotification() {

}
func (t TestManager) ExecuteNotification(task storage.Task, action string) {

}
func (t TestManager) CloseReceivingChannel() {

}
