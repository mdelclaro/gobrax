package driver

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v3"
	"github.com/mdelclaro/gobrax/src/api/helpers"
	database "github.com/mdelclaro/gobrax/src/db"
	"github.com/mdelclaro/gobrax/src/repository/entities"
	"github.com/mdelclaro/gobrax/src/shared"
)

func SetupDriverRoutes(router fiber.Router) {
	driver := router.Group("/driver")
	driver.Get("/:id", GetDriverByID)
	driver.Get("/", GetAllDrivers)
	driver.Post("/", AddDriver)
	driver.Put("/", UpdateDriver)
	driver.Delete("/:id", DeleteDriver)
}

func GetAllDrivers(c fiber.Ctx) error {
	drivers := []entities.Driver{}

	if err := shared.InitRepo(database.DB.Db).FindAll(&drivers); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(helpers.BuildError(err))
	}

	if len(drivers) == 0 {
		return c.Status(http.StatusNoContent).JSON(helpers.ParseResultToMap(""))
	}

	return c.Status(http.StatusOK).JSON(helpers.ParseResultToMap(drivers))
}

func GetDriverByID(c fiber.Ctx) error {
	driver := entities.Driver{}

	id := c.Params("id")
	parsedId, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(helpers.BuildError(fmt.Errorf("invalid id provided: %s", err.Error())))
	}

	if err := shared.InitRepo(database.DB.Db).FindById(&driver, int32(parsedId)); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(helpers.BuildError(err))
	}

	if driver.ID == 0 {
		return c.Status(http.StatusNoContent).JSON(helpers.ParseResultToMap(""))
	}

	return c.Status(http.StatusOK).JSON(helpers.ParseResultToMap(driver))
}

func AddDriver(c fiber.Ctx) error {
	driver := entities.Driver{}

	if err := json.Unmarshal(c.Body(), &driver); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(helpers.BuildError(err))
	}

	if err := helpers.Validator.Struct(driver); err != nil {
		required := ""

		var valErrs validator.ValidationErrors
		if errors.As(err, &valErrs) {

			for i, err := range valErrs {
				field := err.Field()
				if i != 0 {
					field = ", " + field
				}

				required = required + field
			}
		}

		return c.Status(http.StatusBadRequest).JSON(helpers.BuildError(fmt.Errorf("missing required field(s): %s", required)))
	}

	if err := shared.InitRepo(database.DB.Db).Create(&driver); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(helpers.BuildError(err))
	}

	return c.Status(http.StatusCreated).JSON(helpers.ParseResultToMap(driver))
}

func UpdateDriver(c fiber.Ctx) error {
	driver := entities.Driver{}

	if err := json.Unmarshal(c.Body(), &driver); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err.Error())
	}

	if driver.ID == 0 {
		return c.Status(http.StatusBadRequest).JSON(helpers.BuildError(fmt.Errorf("id is required")))
	}

	if err := shared.InitRepo(database.DB.Db).Update(&driver); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(helpers.BuildError(err))
	}

	return c.Status(http.StatusOK).JSON(helpers.ParseResultToMap(driver))
}

func DeleteDriver(c fiber.Ctx) error {
	id := c.Params("id")
	parsedId, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(helpers.BuildError(fmt.Errorf("invalid id provided: %s", err.Error())))
	}

	if err := shared.InitRepo(database.DB.Db).Delete(entities.Driver{}, int32(parsedId)); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(helpers.BuildError(err))
	}

	return c.Status(http.StatusOK).JSON(helpers.ParseResultToMap(""))
}
