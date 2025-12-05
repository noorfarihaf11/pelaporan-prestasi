package routes

import (
	"database/sql"

	"pelaporan-prestasi/app/service"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func StudentRoutes(api fiber.Router, db *sql.DB,  mongoDB *mongo.Database) {

	api.Get("/students", func(c *fiber.Ctx) error {
		return service.GetAllStudentService(c, db)
	})
	api.Get("/students/:id", func(c *fiber.Ctx) error {
		return service.GetStudentByIDService(c, db)
	})
	api.Get("/students/:id/achievements", func(c *fiber.Ctx) error {
		return service.GetAchievementsByStudentIDService(c, mongoDB, db)
	})
	api.Put("/students/:id/advisor", func(c *fiber.Ctx) error {
		return service.UpdateStudentAdvisorService(c, db)
	})

}
