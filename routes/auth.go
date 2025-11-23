package routes

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"pelaporan-prestasi/app/service"
)

func AuthRoutes(api fiber.Router, db *sql.DB) {

	api.Post("/auth/login", func(c *fiber.Ctx) error {
		return service.LoginService(c, db)
	})
	api.Post("/auth/register", func(c *fiber.Ctx) error {
		return service.RegisterService(c, db)
	})
	api.Post("/auth/refresh", func(c *fiber.Ctx) error {
		return service.RefreshTokenService(c, db)
	})

}
