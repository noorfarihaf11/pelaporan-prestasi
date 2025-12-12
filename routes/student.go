package routes

import (
	"database/sql"

	"pelaporan-prestasi/app/service"
	"pelaporan-prestasi/middleware"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func StudentRoutes(api fiber.Router, db *sql.DB,  mongoDB *mongo.Database) {

	api.Get("/students", middleware.RBAC("student:read", db), func(c *fiber.Ctx) error {
		return service.GetAllStudentService(c, db)
	})
	api.Get("/students/:id", middleware.RBAC("student:read", db), func(c *fiber.Ctx) error {
		return service.GetStudentByIDService(c, db)
	})
	api.Get("/students/:id/achievements", middleware.RBAC("student:achievement", db), func(c *fiber.Ctx) error {
		return service.GetAchievementsByStudentIDService(c, mongoDB, db)
	})
	api.Put("/students/:id/advisor", middleware.RBAC("student:advisor", db), func(c *fiber.Ctx) error {
		return service.UpdateStudentAdvisorService(c, db)
	})

}
