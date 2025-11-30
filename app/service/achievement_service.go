package service

import (
	"database/sql"
	"encoding/json"
	"pelaporan-prestasi/app/model"
	"pelaporan-prestasi/app/repository"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateAchievementService(c *fiber.Ctx, mongoDB *mongo.Database, db *sql.DB) error {
    form, err := c.MultipartForm()
    if err != nil {
        return c.Status(400).JSON(fiber.Map{
            "status":  "error",
            "message": "invalid_form_data",
        })
    }

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

    tags := []string{}
    if len(form.Value["tags"]) > 0 {
        json.Unmarshal([]byte(form.Value["tags"][0]), &tags)
    }

    var details map[string]interface{}
    if len(form.Value["details"]) > 0 {
        json.Unmarshal([]byte(form.Value["details"][0]), &details)
    }

    files := form.File["attachments"]
    var attachments []model.Attachment

    for _, file := range files {
        path := "uploads/" + file.Filename

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

     result, err := repository.CreateAchievement(mongoDB, &ach)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "status":  "error",
            "message": "failed_create_achievement_mongo",
        })
    }

    studentUUID, _ := uuid.Parse(ach.StudentID)

    err = repository.CreateAchievementReference(db, studentUUID, result.ID.Hex())
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "status":  "error",
            "message": "failed_create_achievement_reference_postgres",
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

func UpdateAchievementService(c *fiber.Ctx, mongoDB *mongo.Database, db *sql.DB) error {
    id := c.Params("id") // mongo achievement id (hex)
    if id == "" {
        return c.Status(400).JSON(fiber.Map{"status": "error", "message": "missing_id"})
    }

    form, err := c.MultipartForm()
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"status": "error", "message": "invalid_form_data"})
    }

    title := form.Value["title"][0]
    description := form.Value["description"][0]
    status := "draft"
    if len(form.Value["status"]) > 0 {
        status = form.Value["status"][0]
    }

    var details map[string]interface{}
    if len(form.Value["details"]) > 0 {
        json.Unmarshal([]byte(form.Value["details"][0]), &details)
    }

    tags := []string{}
    if len(form.Value["tags"]) > 0 {
        json.Unmarshal([]byte(form.Value["tags"][0]), &tags)
    }

    files := form.File["attachments"]
    var attachments []model.Attachment

    for _, file := range files {
        path := "uploads/" + file.Filename
        c.SaveFile(file, path)

        attachments = append(attachments, model.Attachment{
            FileName:   file.Filename,
            FileUrl:    path,
            FileType:   file.Header.Get("Content-Type"),
            UploadedAt: time.Now(),
        })
    }

    points := 0
    if len(form.Value["points"]) > 0 {
        points, _ = strconv.Atoi(form.Value["points"][0])
    }

    updated, err := repository.UpdateAchievement(mongoDB, id, bson.M{
        "title":        title,
        "description":  description,
        "details":      details,
        "tags":         tags,
        "attachments":  attachments,
        "status":       status,
        "points":       points,
        "updated_at":   time.Now(),
    })

    if err != nil {
        return c.Status(500).JSON(fiber.Map{"status": "error", "message": "failed_update_mongo"})
    }

    err = repository.UpdateAchievementReference(db, id, status)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"status": "error", "message": "failed_update_reference"})
    }

    return c.Status(200).JSON(fiber.Map{
        "status":  "success",
        "message": "achievement_updated_successfully",
        "data":    updated,
    })
}
func SoftDeleteAchievementService(c *fiber.Ctx, mongoDB *mongo.Database, db *sql.DB) error {
    id := c.Params("id")

    // Soft delete Mongo
    err := repository.SoftDeleteAchievement(mongoDB, id)
    if err != nil {
        return c.Status(400).JSON(fiber.Map{
            "status":  "error",
            "message": err.Error(),
        })
    }

    // Update PostgreSQL reference
    err = repository.UpdateAchievementReference(db, id, "deleted")
    if err != nil {
        msg := err.Error()

        if msg == "reference_not_found" {
            return c.Status(404).JSON(fiber.Map{
                "status":  "error",
                "message": "achievement_reference_not_found",
            })
        }

        return c.Status(500).JSON(fiber.Map{
            "status":  "error",
            "message": "failed_soft_delete_reference",
            "detail":  msg,
        })
    }

    return c.Status(200).JSON(fiber.Map{
        "status":  "success",
        "message": "achievement_deleted_successfully",
    })
}

func SubmitAchievementService(c *fiber.Ctx, mongoDB *mongo.Database, db *sql.DB) error {
    id := c.Params("id")

    err := repository.SubmitAchievement(mongoDB, id)
    if err != nil {
        return c.Status(400).JSON(fiber.Map{
            "status":  "error",
            "message": err.Error(),
        })
    }

    err = repository.UpdateAchievementReference(db, id, "submitted")
    if err != nil {
        msg := err.Error()

        if msg == "reference_not_found" {
            return c.Status(404).JSON(fiber.Map{
                "status":  "error",
                "message": "achievement_reference_not_found",
            })
        }

        return c.Status(500).JSON(fiber.Map{
            "status":  "error",
            "message": "failed submit",
            "detail":  msg,
        })
    }

    return c.Status(200).JSON(fiber.Map{
        "status":  "success",
        "message": "achievement submitted successfully",
    })
}
