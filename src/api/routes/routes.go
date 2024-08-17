package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/mdelclaro/gobrax/src/api/handlers/driver"
	"github.com/mdelclaro/gobrax/src/api/handlers/truck"
)

func SetUpRoutes(app *fiber.App) {
	api := app.Group("/api", logger.New())

	driver.SetupDriverRoutes(api)
	truck.SetupTruckRoutes(api)
}
