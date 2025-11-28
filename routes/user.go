package routes

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"pelaporan-prestasi/app/service"
	"pelaporan-prestasi/middleware"
)

func UserRoutes(api fiber.Router, db *sql.DB) {

	api.Get("/users", middleware.RBAC("user:manage", db), func(c *fiber.Ctx) error {
		return service.GetAllUserService(c, db)
	})

	api.Get("/users/:id", func(c *fiber.Ctx) error {
		return service.GetUserByIDService(c, db)
	})
	api.Post("/users", func(c *fiber.Ctx) error {
		return service.CreateUserService(c, db)
	})
	api.Put("/users/:id", func(c *fiber.Ctx) error {
		return service.UpdateUserService(c, db)
	})
	api.Delete("/users/:id", func(c *fiber.Ctx) error {
		return service.DeleteUserService(c, db)
	})
	api.Put("/users/:id/role", func(c *fiber.Ctx) error {
		return service.UpdateUserRoleService(c, db)
	})

}
