package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/mdelclaro/gobrax/src/api/handlers"
)

func SetupDriverRoutes(router fiber.Router) {
	driver := router.Group("/driver")
	driver.Get("/:id", handlers.GetDriverByID)
	driver.Get("/", handlers.GetAllDrivers)
	driver.Post("/", handlers.AddDriver)
	driver.Put("/", handlers.UpdateDriver)
	driver.Delete("/:id", handlers.DeleteDriver)
}
