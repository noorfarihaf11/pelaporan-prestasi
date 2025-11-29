package routes

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func Routes(app *fiber.App, db *sql.DB, mongoDB *mongo.Database) {
	api := app.Group("/api/v1/")

	AuthRoutes(api, db)
	UserRoutes(api, db)
	StudentRoutes(api, db)
	LecturerRoutes(api, db)
	AchievementRoutes(api, db, mongoDB)
}

