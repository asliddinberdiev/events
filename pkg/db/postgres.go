package db

import (
	"fmt"
	"log"
	"time"

	"github.com/asliddinberdiev/events/conf"
	"github.com/jmoiron/sqlx"
)

var (
	DB              *sqlx.DB
	maxOpenConns    = 20
	maxIdleConns    = 20
	maxConnLifeTime = 10
)

func ConnectPSQL() {
	dbUrl := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", conf.Envs.Postgres.Host, conf.Envs.Postgres.Port, conf.Envs.Postgres.User, conf.Envs.Postgres.Name, conf.Envs.Postgres.Password, conf.Envs.Postgres.SSLMode)
	db, err := sqlx.Connect("postgres", dbUrl)
	if err != nil {
		log.Fatalf("Error: Failed to connect to the database: %v", err)
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(time.Duration(maxConnLifeTime) * time.Second)

	if err := db.Ping(); err != nil {
		log.Fatal("Error: Ping request")
	}

	DB = db
}

func DisConnectPSQL() {
	DB.Close()
}
