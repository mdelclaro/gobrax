package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/mdelclaro/gobrax/src/config"
	database "github.com/mdelclaro/gobrax/src/db"
	"github.com/mdelclaro/gobrax/src/shared"
	"github.com/mdelclaro/gobrax/src/utils"
)

func main() {
	database.StartDb()
	shared.InitRepo(database.DB.Db)
	app := utils.SetupApp()

	app.Use(cors.New())
	app.Use(func(c fiber.Ctx) error {
		return fiber.ErrNotFound
	})

	log.Fatal(app.Listen(config.GetEnv("APP_PORT")))
}
