package routes

import (
	"database/sql"
	"pelaporan-prestasi/app/service"
	"pelaporan-prestasi/middleware"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func AchievementRoutes(api fiber.Router, db *sql.DB, mongoDB *mongo.Database) {

    api.Post("/achievements", middleware.RBAC("achievement:create", db), func(c *fiber.Ctx) error {
        return service.CreateAchievementService(c, mongoDB, db)
    })

	api.Get("/achievements", func(c *fiber.Ctx) error {
		return service.GetAllAchievementsService(c, mongoDB, db)
	})

	api.Get("/achievements/:id", func(c *fiber.Ctx) error {
		return service.GetAchievementByIDService(c, mongoDB, db)
	})

	api.Put("/achievements/:id", func(c *fiber.Ctx) error {
        return service.UpdateAchievementService(c, mongoDB, db)
    })

	api.Delete("/achievements/delete/:id", func(c *fiber.Ctx) error {
        return service.SoftDeleteAchievementService(c, mongoDB, db)
    })

	api.Post("/achievements/:id/submit", middleware.RBAC("achievement:submit", db), func(c *fiber.Ctx) error {
        return service.SubmitAchievementService(c, mongoDB, db)
    })

	api.Post("/achievements/:id/verify", middleware.RBAC("achievement:verify", db), func(c *fiber.Ctx) error {
        return service.VerifyAchievementService(c, mongoDB, db)
    })

	api.Post("/achievements/:id/reject", middleware.RBAC("achievement:reject", db), func(c *fiber.Ctx) error {
        return service.RejectAchievementService(c, mongoDB, db)
    })

	api.Get("/reports/statistics", func(c *fiber.Ctx) error {
		return service.GetStatisticsService(c, mongoDB)
	})

}
