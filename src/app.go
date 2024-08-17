package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/mdelclaro/gobrax/src/api/routes"
	database "github.com/mdelclaro/gobrax/src/db"
	"github.com/mdelclaro/gobrax/src/shared"
)

func main() {
	database.StartDb()
	shared.InitRepo()
	app := fiber.New()

	routes.SetUpRoutes(app)

	app.Use(cors.New())
	app.Use(func(c fiber.Ctx) error {
		return fiber.ErrNotFound
	})

	log.Fatal(app.Listen(":3000"))
}
