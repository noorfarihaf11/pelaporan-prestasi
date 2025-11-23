package config

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	// _"pelaporan-prestasi/middleware"
	"pelaporan-prestasi/routes"
)

func NewApp(db *sql.DB) *fiber.App {
	app := fiber.New()

	routes.AuthRoutes(app, db)

	return app
}
