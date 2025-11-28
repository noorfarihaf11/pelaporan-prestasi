package routes

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"pelaporan-prestasi/app/service"
)

func LecturerRoutes(api fiber.Router, db *sql.DB) {
	api.Get("/lecturers", func(c *fiber.Ctx) error {
		return service.GetAllLecturersService(c, db)
	})
}
