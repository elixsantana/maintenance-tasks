package storage

import "os"

type MysqlConfig struct {
	User     string
	Password string
	Host     string
	Port     int
}

func LoadMysqlConfig() (*MysqlConfig, error) {
	actualHost := "localhost"
	host, exists := os.LookupEnv("LOCALHOST")
	if exists {
		actualHost = host
	}

	user := "root"
	password := "test"
	port := 3306

	return &MysqlConfig{
		User:     user,
		Password: password,
		Host:     actualHost,
		Port:     port,
	}, nil
}
