package routes

import (
	"database/sql"
	"pelaporan-prestasi/app/service"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func AchievementRoutes(api fiber.Router, db *sql.DB, mongoDB *mongo.Database) {

    api.Post("/achievements", func(c *fiber.Ctx) error {
        return service.CreateAchievementService(c, mongoDB, db)
    })

	api.Get("/achievements", func(c *fiber.Ctx) error {
		return service.GetAllAchievementsService(c, mongoDB)
	})

	api.Get("/achievements/:id", func(c *fiber.Ctx) error {
		return service.GetAchievementByIDService(c, mongoDB)
	})

	api.Put("/achievements/:id", func(c *fiber.Ctx) error {
        return service.UpdateAchievementService(c, mongoDB, db)
    })

	api.Put("/achievements/delete/:id", func(c *fiber.Ctx) error {
        return service.SoftDeleteAchievementService(c, mongoDB, db)
    })
}
