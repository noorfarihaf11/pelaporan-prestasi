package routes

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"pelaporan-prestasi/app/service"
)

func UserRoutes(api fiber.Router, db *sql.DB) {

	api.Get("/users", func(c *fiber.Ctx) error {
		return service.GetAllUserService(c, db)
	})
	api.Get("/users/:id", func(c *fiber.Ctx) error {
		return service.GetUserByIDService(c, db)
	})
	api.Post("/users", func(c *fiber.Ctx) error {
		return service.CreateUserService(c, db)
	})

}
