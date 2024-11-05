package conf

import (
	"log"
	"os"
	"strconv"
)

type EnvConf struct {
	App      App
	Postgres Postgres
}

type App struct {
	JWTSecret              string
	JWTExpirationInSeconds int64
}

type Postgres struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

var Envs EnvConf

func LoadConf() {
	Envs.App.JWTSecret = getEnv("JWT_SECRET", "secret_key_not_found!")
	jwtExp, err := strconv.ParseInt(getEnv("JWT_EXP", "604800"), 10, 64)
	if err != nil {
		log.Printf("Warning: Invalid JWT_EXP, using default value: 604800")
		jwtExp = 604800
	}
	Envs.App.JWTExpirationInSeconds = jwtExp

	Envs.Postgres.Host = getEnv("POSTGRES_HOST", "localhost")
	Envs.Postgres.Port = getEnv("POSTGRES_PORT", "5432")
	Envs.Postgres.User = getEnv("POSTGRES_USER", "user")
	Envs.Postgres.Password = getEnv("POSTGRES_PASSWORD", "password")
	Envs.Postgres.Name = getEnv("POSTGRES_DB", "name")
	Envs.Postgres.SSLMode = getEnv("POSTGRES_SSLMODE", "disable")
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}
