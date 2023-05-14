package manager

import (
	"fmt"
	"log"
	database "maintenance-tasks/storage"
	"strconv"
)

type Manager struct {
	databaseMetadata *database.MysqlMetadata
}

func Create() *Manager {
	config, err := database.LoadMysqlConfig()
	if err != nil {
		log.Fatal(err)
	}

	mysqlCreate := database.Create(config)

	return &Manager{
		databaseMetadata: mysqlCreate,
	}
}

func (m *Manager) Start() {
	err := m.databaseMetadata.Connect()
	if err != nil {
		log.Fatal(err)
	}

	err = m.databaseMetadata.CreateTaskTable()
	if err != nil {
		log.Fatal(err)
	}

}

func (m *Manager) Stop() {
	m.databaseMetadata.Close()
}

func (m *Manager) GetAllTasks() ([]database.Task, error) {
	tasks, err := m.databaseMetadata.GetAllTasks()
	if err != nil {
		return []database.Task{}, err
	}

	return tasks, err
}

func (m *Manager) CreateTask(summary string, techId string, role string) error {
	id, err := strconv.Atoi(techId)
	if id < 0 {
		return fmt.Errorf("not valid ID")
	}
	if err != nil {
		return err
	}

	err = m.databaseMetadata.CreateTask(summary, id, role)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) UpdateTask(task database.Task) (database.Task, error) {
	task, err := m.databaseMetadata.UpdateTask(task)
	if err != nil {
		return database.Task{}, err
	}

	return task, err
}

func (m *Manager) GetTask(task_id int, tech_id int) (database.Task, error) {
	task, err := m.databaseMetadata.GetTask(task_id, tech_id)
	if err != nil {
		return database.Task{}, err
	}

	return task, err

}

func (m *Manager) DeleteTask(taskID int) error {
	err := m.databaseMetadata.DeleteTask(taskID)
	if err != nil {
		return err
	}

	return nil
}
