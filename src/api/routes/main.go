package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

func SetUpRoutes(app *fiber.App) {
	api := app.Group("/api", logger.New())

	SetupDriverRoutes(api)
	SetupTruckRoutes(api)
}
