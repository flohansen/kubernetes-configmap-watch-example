package postgres

import "fmt"

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

func (cfg *Config) Dsn() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
}
