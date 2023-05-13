package manager

import (
	"log"
	database "maintenance-tasks/storage"
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
