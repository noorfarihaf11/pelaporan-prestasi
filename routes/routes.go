package routes

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

func Routes(app *fiber.App, db *sql.DB) {
	api := app.Group("/api/v1/")

	AuthRoutes(api, db)
	UserRoutes(api, db)
	StudentRoutes(api, db)
}
