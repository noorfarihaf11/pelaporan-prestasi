package routes

import (
	"pelaporan-prestasi/app/service"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func AchievementRoutes(api fiber.Router, mongoDB *mongo.Database) {
	api.Post("/achievements", func(c *fiber.Ctx) error {
		return service.CreateAchievementService(c, mongoDB)
	})

	api.Get("/achievements", func(c *fiber.Ctx) error {
		return service.GetAllAchievementsService(c, mongoDB)
	})

	api.Get("/achievements/:id", func(c *fiber.Ctx) error {
		return service.GetAchievementByIDService(c, mongoDB)
	})
}
