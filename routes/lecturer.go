package routes

import (
	"database/sql"

	"pelaporan-prestasi/app/service"
	"pelaporan-prestasi/middleware"

	"github.com/gofiber/fiber/v2"
)

func LecturerRoutes(api fiber.Router, db *sql.DB) {
	api.Get("/lecturers", middleware.RBAC("lecturer:read", db), func(c *fiber.Ctx) error {
		return service.GetAllLecturersService(c, db)
	})
	api.Get("/lecturers/:id/advisees", middleware.RBAC("lecturer:read", db), func(c *fiber.Ctx) error {
		return service.GetAdviseesByLecturerService(c, db)
	})
}
