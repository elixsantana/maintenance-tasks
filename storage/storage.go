package storage

import (
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

type MysqlMetadata struct {
	config *MysqlConfig
	uri    string
	mutex  sync.Mutex
	db     *sql.DB
}

func Create(mysqlConfig *MysqlConfig) *MysqlMetadata {
	return &MysqlMetadata{
		config: mysqlConfig,
		uri: fmt.Sprintf("%v:%v@tcp(%v:%v)/",
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
		roles ENUM('technician', 'manager') NOT NULL DEFAULT 'technician'
    )`)
	if err != nil {
		return err
	}

	return nil
}
