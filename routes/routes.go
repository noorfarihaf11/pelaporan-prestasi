package routes

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

func Routes(app *fiber.App, db *sql.DB, mongoDB *mongo.Database) {
	api := app.Group("/api/v1/")
	app.Get("/swagger/*", fiberSwagger.WrapHandler)
	AuthRoutes(api, db)
	UserRoutes(api, db)
	StudentRoutes(api, db, mongoDB)
	LecturerRoutes(api, db)
	AchievementRoutes(api, db, mongoDB)
	StatisticRoutes(api, db, mongoDB)
}

