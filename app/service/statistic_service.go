package service

import (
	"pelaporan-prestasi/app/repository"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetStatisticsService(c *fiber.Ctx, db *mongo.Database) error {
	stats, err := repository.GetAchievementStatistics(db)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed fetch statistic",
			"detail":  err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Success get statistic",
		"data":    stats,
	})
}
func GetStudentStatisticsService(c *fiber.Ctx, db *mongo.Database) error {
	studentID := c.Params("id")

	stats, err := repository.GetStudentStatistics(db, studentID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed fetch students statistic",
			"detail":  err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Success get student statistic",
		"data":    stats,
	})
}
