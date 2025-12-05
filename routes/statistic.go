package routes

import (
	"database/sql"
	"pelaporan-prestasi/app/service"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func StatisticRoutes(api fiber.Router, db *sql.DB, mongoDB *mongo.Database) {

	api.Get("/reports/statistics", func(c *fiber.Ctx) error {
		return service.GetStatisticsService(c, mongoDB)
	})

	api.Get("/reports/student/:id", func(c *fiber.Ctx) error {
		return service.GetStudentStatisticsService(c, mongoDB)
	})

}
