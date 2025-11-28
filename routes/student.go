package routes

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"pelaporan-prestasi/app/service"
)

func StudentRoutes(api fiber.Router, db *sql.DB) {

	api.Get("/students", func(c *fiber.Ctx) error {
		return service.GetAllStudentService(c, db)
	})
	api.Get("/students/:id", func(c *fiber.Ctx) error {
		return service.GetStudentByIDService(c, db)
	})

}
