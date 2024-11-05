package main

import (
	"log"

	"github.com/asliddinberdiev/events/conf"
	"github.com/asliddinberdiev/events/internal/http"
	"github.com/asliddinberdiev/events/pkg/db"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error: loading .env file")
	}
	conf.LoadConf()

	db.ConnectPSQL()
	defer db.DisConnectPSQL()

	app := http.NewRouter()

	if err := app.Listen(":3000"); err != nil {
		log.Fatal(err)
	}
}
