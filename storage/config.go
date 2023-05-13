package storage

type MysqlConfig struct {
	User     string
	Password string
	Host     string
	Port     int
}

func LoadMysqlConfig() (*MysqlConfig, error) {
	user := "root"
	password := "test"
	host := "localhost"
	port := 3306

	return &MysqlConfig{
		User:     user,
		Password: password,
		Host:     host,
		Port:     port,
	}, nil
}
