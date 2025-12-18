package service

import (
	"pelaporan-prestasi/app/repository"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetStatistics godoc
// @Summary      Get achievement statistics
// @Description  Mengambil statistik pencapaian secara keseluruhan
// @Tags         Statistic
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} model.BaseResponse{data=map[string]interface{}}
// @Failure      500 {object} model.BaseResponse
// @Router       /api/v1/statistics [get]
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

// GetStudentStatistics godoc
// @Summary      Get student achievement statistics
// @Description  Mengambil statistik pencapaian berdasarkan student ID
// @Tags         Statistic
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Student ID"
// @Success      200 {object} model.BaseResponse{data=map[string]interface{}}
// @Failure      500 {object} model.BaseResponse
// @Router       /api/v1/statistics/students/{id} [get]
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
