package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/mdelclaro/gobrax/src/api/handlers"
)

func SetupTruckRoutes(router fiber.Router) {
	truck := router.Group("/truck")
	truck.Get("/:id", handlers.GetTruckByID)
	truck.Get("/", handlers.GetAllTrucks)
	truck.Post("/", handlers.AddTruck)
	truck.Put("/", handlers.UpdateTruck)
	truck.Delete("/:id", handlers.DeleteTruck)
	truck.Post("/update-driver/:id", handlers.UpdateTruckDriver)
}
