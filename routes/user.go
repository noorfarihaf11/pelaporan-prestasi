package routes

import (
	"database/sql"

	"pelaporan-prestasi/app/service"
	"pelaporan-prestasi/middleware"

	"github.com/gofiber/fiber/v2"
)

func UserRoutes(api fiber.Router, db *sql.DB) {

	api.Get("/users", middleware.RBAC("user:manage", db), func(c *fiber.Ctx) error {
		return service.GetAllUserService(c, db)
	})

	api.Get("/users/:id", middleware.RBAC("user:manage", db), func(c *fiber.Ctx) error {
		return service.GetUserByIDService(c, db)
	})
	api.Post("/users", middleware.RBAC("user:manage", db), func(c *fiber.Ctx) error {
		return service.CreateUserService(c, db)
	})
	api.Put("/users/:id", middleware.RBAC("user:manage", db), func(c *fiber.Ctx) error {
		return service.UpdateUserService(c, db)
	})
	api.Delete("/users/:id", middleware.RBAC("user:manage", db), func(c *fiber.Ctx) error {
		return service.DeleteUserService(c, db)
	})
	api.Put("/users/:id/role", middleware.RBAC("user:manage", db), func(c *fiber.Ctx) error {
		return service.UpdateUserRoleService(c, db)
	})

}
