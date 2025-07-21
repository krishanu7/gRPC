package database

import (
	"database/sql"
	"fmt"
	"github.com/krishanu7/grpc/internal/config"
	_ "github.com/lib/pq"
)

func NewDB(cfg *config.Config) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName)
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil , err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, err
}