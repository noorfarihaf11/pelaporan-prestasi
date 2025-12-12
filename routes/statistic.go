package routes

import (
	"database/sql"
	"pelaporan-prestasi/app/service"
	"pelaporan-prestasi/middleware"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func StatisticRoutes(api fiber.Router, db *sql.DB, mongoDB *mongo.Database) {

	api.Get("/reports/statistics", middleware.RBAC("report:statistic", db), func(c *fiber.Ctx) error {
		return service.GetStatisticsService(c, mongoDB)
	})

	api.Get("/reports/student/:id", middleware.RBAC("report:statistic", db), func(c *fiber.Ctx) error {
		return service.GetStudentStatisticsService(c, mongoDB)
	})

}
