package main

import (
	"log"
	"os"
	"pelaporan-prestasi/config"
	"pelaporan-prestasi/database"
	"pelaporan-prestasi/routes"
	"time"
)

func main() {
	config.LoadEnv()
	db := database.ConnectDB()
	app := config.NewApp(db)
	port := os.Getenv("APP_PORT")
	routes.Routes(app, db)
	if port == "" {
		port = "3000"
	}

	loc, _ := time.LoadLocation("Asia/Jakarta")
	time.Local = loc

	log.Fatal(app.Listen(":" + port))
}
