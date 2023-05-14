package storage

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type MysqlMetadata struct {
	config *MysqlConfig
	uri    string
	mutex  sync.Mutex
	db     *sql.DB
}

type Task struct {
	ID             uint      `json:"id"`
	Summary        string    `json:"summary"`
	Performed_date time.Time `json:"performed_date"`
	TechnicianID   int       `json:"technician_id"`
	Role           string    `json:"role"`
}

func Create(mysqlConfig *MysqlConfig) *MysqlMetadata {
	return &MysqlMetadata{
		config: mysqlConfig,
		uri: fmt.Sprintf("%v:%v@tcp(%v:%v)/?parseTime=true",
			mysqlConfig.User,
			mysqlConfig.Password,
			mysqlConfig.Host,
			mysqlConfig.Port,
		),
	}
}

func (m *MysqlMetadata) Connect() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.db != nil {
		return fmt.Errorf("connection is open")
	}

	var err error
	m.db, err = sql.Open("mysql", m.uri)

	if err != nil {
		return fmt.Errorf("failed connection: %v", err)
	}

	return nil
}

func (m *MysqlMetadata) Close() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.db != nil {
		err := m.db.Close()
		m.db = nil
		return err
	}
	return nil
}

func (m *MysqlMetadata) CreateTaskTable() error {
	_, err := m.db.Exec("CREATE DATABASE IF NOT EXISTS maintenance")
	if err != nil {
		return err
	}

	_, err = m.db.Exec("use maintenance")
	if err != nil {
		return err
	}

	_, err = m.db.Exec(`CREATE TABLE IF NOT EXISTS tasks (
        id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
        summary VARCHAR(2500) NOT NULL,
        performed_date DATE NOT NULL,
        technician_id INT NOT NULL,
		role ENUM('technician', 'manager') NOT NULL DEFAULT 'technician'
    )`)
	if err != nil {
		return err
	}

	return nil
}

func (m *MysqlMetadata) GetAllTasks() ([]Task, error) {
	rows, err := m.db.Query("SELECT * FROM tasks")
	if err != nil {
		return []Task{}, err
	}
	defer rows.Close()

	var tasks []Task

	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Summary, &task.Performed_date, &task.TechnicianID, &task.Role)
		if err != nil {
			return []Task{}, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (m *MysqlMetadata) CreateTask(summary string, tech_id int, role string) error {
	now := time.Now()
	role = strings.ToLower(role)

	stmt, err := m.db.Prepare("INSERT INTO tasks(summary, performed_date, technician_id, role) VALUES (?, ?, ?, ?)")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	_, err = stmt.Exec(summary, now, tech_id, role)
	if err != nil {
		return err
	}

	return nil

}

func (m *MysqlMetadata) GetTask(id int, techId int) (Task, error) {
	var task Task
	query := "SELECT * FROM tasks where id= ? AND technician_id=?"

	stmt, err := m.db.Prepare(query)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	row := stmt.QueryRow(id, techId)
	switch err := row.Scan(&task.ID, &task.Summary, &task.Performed_date, &task.TechnicianID, &task.Role); err {
	case sql.ErrNoRows:
		//fmt.Println("No rows were returned!")
		return Task{}, err
	case nil:
		return task, nil
	default:
		panic(err)
	}

}

func (m *MysqlMetadata) UpdateTask(task Task) (Task, error) {
	query := "UPDATE tasks SET summary=?, performed_date=?, technician_id=?, role=? WHERE id=?"
	stmt, err := m.db.Prepare(query)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	result, err := stmt.Exec(task.Summary, task.Performed_date, task.TechnicianID, task.Role, task.ID)
	if err != nil {
		fmt.Println(err.Error())
		return task, fmt.Errorf("%d", http.StatusInternalServerError)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Println(err.Error())
		return task, fmt.Errorf("%d", http.StatusInternalServerError)
	}

	if rowsAffected == 0 {
		return task, fmt.Errorf("no changes")

	}

	return task, nil
}
