package service

import (
	"encoding/json"
	"pelaporan-prestasi/app/model"
	"pelaporan-prestasi/app/repository"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateAchievementService(c *fiber.Ctx, mongoDB *mongo.Database) error {
    // Ambil semua form-data
    form, err := c.MultipartForm()
    if err != nil {
        return c.Status(400).JSON(fiber.Map{
            "status":  "error",
            "message": "invalid_form_data",
        })
    }

    // Extract fields
    studentID := form.Value["student_id"][0]
    achievementType := form.Value["achievement_type"][0]
    title := form.Value["title"][0]
    description := form.Value["description"][0]
    points := 0

	status := "draft"
	if len(form.Value["status"]) > 0 {
		status = form.Value["status"][0]
	}


    if len(form.Value["points"]) > 0 {
        points, _ = strconv.Atoi(form.Value["points"][0])
    }

    // Parse tags
    tags := []string{}
    if len(form.Value["tags"]) > 0 {
        json.Unmarshal([]byte(form.Value["tags"][0]), &tags)
    }

    // Parse details (dinamis)
    var details map[string]interface{}
    if len(form.Value["details"]) > 0 {
        json.Unmarshal([]byte(form.Value["details"][0]), &details)
    }

    // Handle attachments
    files := form.File["attachments"]
    var attachments []model.Attachment

    for _, file := range files {
        // Simpan file ke folder lokal (atau ke cloud)
        path := "uploads/" + file.Filename

        // Save file
        if err := c.SaveFile(file, path); err != nil {
            return c.Status(500).JSON(fiber.Map{
                "status":  "error",
                "message": "failed_save_file",
            })
        }

        attachments = append(attachments, model.Attachment{
            FileName: file.Filename,
            FileUrl:  path,
            FileType: file.Header.Get("Content-Type"),
            UploadedAt: time.Now(),
        })
    }

    // Bentuk struct Achievement
    ach := model.Achievement{
        StudentID: studentID,
        AchievementType: achievementType,
        Title: title,
        Description: description,
        Details: details,
        Tags: tags,
        Points: points,
        Attachments: attachments,
		Status: status,
    }

    // Insert ke MongoDB
    result, err := repository.CreateAchievement(mongoDB, &ach)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "status":  "error",
            "message": "failed_create_achievement",
        })
    }

    return c.Status(201).JSON(fiber.Map{
        "status":  "success",
        "message": "achievement_created_successfully",
        "data": fiber.Map{
            "achievement": result,
        },
    })
}

func GetAllAchievementsService(c *fiber.Ctx, mongoDB *mongo.Database) error {
	list, err := repository.GetAllAchievements(mongoDB)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "failed_fetch_achievements",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "success_get_all_achievements",
		"data": fiber.Map{
			"achievements": list,
		},
	})
}
func GetAchievementByIDService(c *fiber.Ctx, mongoDB *mongo.Database) error {
	id := c.Params("id")

	ach, err := repository.GetAchievementByID(mongoDB, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "failed_fetch_achievement",
		})
	}

	if ach == nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "achievement_not_found",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "success_get_achievement",
		"data": fiber.Map{
			"achievement": ach,
		},
	})
}

