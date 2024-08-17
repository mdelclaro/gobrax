package utils

import (
	"github.com/gofiber/fiber/v3"
	"github.com/mdelclaro/gobrax/src/api/routes"
)

func SetupApp() *fiber.App {
	app := fiber.New()

	routes.SetUpRoutes(app)

	return app
}
