package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/mdelclaro/gobrax/src/api/handlers"
)

func SetupTruckRoutes(router fiber.Router) {
	truck := router.Group("/truck")
	truck.Get("/", handlers.GetAllTrucks)
}
