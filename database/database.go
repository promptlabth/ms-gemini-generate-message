package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Config struct {
	UserName     string
	Password     string
	Host         string
	Port         string
	DatabaseName string
}

func Open(cfg Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.Port, cfg.UserName, cfg.Password, cfg.DatabaseName))
	if err != nil {
		return nil, err
	}
	return db, nil
}
