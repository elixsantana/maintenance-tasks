package manager

import (
	"fmt"
	"log"
	database "maintenance-tasks/storage"
	"time"
)

type TaskAction struct {
	task   database.Task
	action string
}

type Manager interface {
	Start()
	Stop()
	GetAllTasks() ([]database.Task, error)
	CreateTask(summary string, techId int, role string, now time.Time) error
	UpdateTask(task database.Task) (database.Task, error)
	GetTask(task_id int, tech_id int, manager bool) (database.Task, error)
	DeleteTask(taskID int) error
	ReceiveNotification()
	ExecuteNotification(task database.Task, action string)
	CloseReceivingChannel()
}

type ManagerImpl struct {
	databaseMetadata *database.MysqlMetadata
	taskCh           chan TaskAction
	doneCh           chan bool
}

func Create() *ManagerImpl {
	config, err := database.LoadMysqlConfig()
	if err != nil {
		log.Fatal(err)
	}

	mysqlCreate := database.Create(config)

	return &ManagerImpl{
		databaseMetadata: mysqlCreate,
		taskCh:           make(chan TaskAction),
		doneCh:           make(chan bool),
	}
}

func (m *ManagerImpl) Start() {
	err := m.databaseMetadata.Connect()
	if err != nil {
		log.Fatal(err)
	}

	err = m.databaseMetadata.CreateTaskTable()
	if err != nil {
		log.Fatal(err)
	}

	m.ReceiveNotification()
}

func (m *ManagerImpl) Stop() {
	m.databaseMetadata.Close()
}

func (m *ManagerImpl) GetAllTasks() ([]database.Task, error) {
	tasks, err := m.databaseMetadata.GetAllTasks()
	if err != nil {
		return []database.Task{}, err
	}

	return tasks, err
}

func (m *ManagerImpl) CreateTask(summary string, techId int, role string, now time.Time) error {
	if techId < 1 {
		return fmt.Errorf("not valid ID")
	}

	err := m.databaseMetadata.CreateTask(summary, techId, role, now)
	if err != nil {
		return err
	}

	return nil
}

func (m *ManagerImpl) UpdateTask(task database.Task) (database.Task, error) {
	task, err := m.databaseMetadata.UpdateTask(task)
	if err != nil {
		return database.Task{}, err
	}

	return task, err
}

func (m *ManagerImpl) GetTask(task_id int, tech_id int, manager bool) (database.Task, error) {
	task, err := m.databaseMetadata.GetTask(task_id, tech_id, manager)
	if err != nil {
		return database.Task{}, err
	}

	return task, err

}

func (m *ManagerImpl) DeleteTask(taskID int) error {
	err := m.databaseMetadata.DeleteTask(taskID)
	if err != nil {
		return err
	}

	return nil
}

func (m *ManagerImpl) ReceiveNotification() {
	var resulTaskAction TaskAction
	go func() {
		for {
			select {
			case resulTaskAction = <-m.taskCh:
				if resulTaskAction.task.ID > 0 {
					fmt.Printf("The tech with id #%d performed the task %s on date %s\n", resulTaskAction.task.TechnicianID, resulTaskAction.action, resulTaskAction.task.Performed_date.Format("2006-01-02 15:04:05"))
				}
			case <-m.doneCh:
				fmt.Println("Stopping notifications")
				return
			}
		}
	}()
}

func (m *ManagerImpl) ExecuteNotification(task database.Task, action string) {
	var result TaskAction
	result.task = task
	result.action = action
	go func() {
		m.taskCh <- result
	}()
}

func (m *ManagerImpl) CloseReceivingChannel() {
	m.doneCh <- true
}
