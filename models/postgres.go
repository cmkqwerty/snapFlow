package models

import (
	"database/sql"
	"fmt"
	"github.com/cmkqwerty/snapFlow/configs"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func Open(config PostgresConfig) (*sql.DB, error) {
	db, err := sql.Open("pgx", config.String())
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}

	return db, nil
}

func DefaultPostgresConfig() PostgresConfig {
	return PostgresConfig{
		Host:     configs.ReadKey("DB_HOST"),
		Port:     configs.ReadKey("DB_PORT"),
		User:     configs.ReadKey("DB_USER"),
		Password: configs.ReadKey("DB_PASSWORD"),
		Database: configs.ReadKey("DB_NAME"),
		SSLMode:  configs.ReadKey("DB_SSLMODE"),
	}
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

func (config PostgresConfig) String() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", config.Host, config.Port, config.User, config.Password, config.Database, config.SSLMode)
}
