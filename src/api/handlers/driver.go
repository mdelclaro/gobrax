package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v3"
	"github.com/mdelclaro/gobrax/src/api/helpers"
	database "github.com/mdelclaro/gobrax/src/db"
	"github.com/mdelclaro/gobrax/src/models"
	"gorm.io/gorm/clause"
)

func GetAllDrivers(c fiber.Ctx) error {
	driver := models.Driver{}

	if tx := database.DB.Db.Find(&driver); tx.Error != nil {
		return c.Status(http.StatusInternalServerError).JSON(helpers.BuildError(tx.Error))
	}

	if driver.ID == 0 {
		return c.Status(http.StatusNoContent).JSON(helpers.ParseResultToMap(""))
	}

	return c.Status(http.StatusOK).JSON(helpers.ParseResultToMap(driver))
}

func GetDriverByID(c fiber.Ctx) error {
	driver := models.Driver{}
	id := c.Params("id")

	if tx := database.DB.Db.First(&driver, id); tx.Error != nil {
		return c.Status(http.StatusInternalServerError).JSON(helpers.BuildError(tx.Error))
	}

	if driver.ID == 0 {
		return c.Status(http.StatusNoContent).JSON(helpers.ParseResultToMap(""))
	}

	return c.Status(http.StatusOK).JSON(helpers.ParseResultToMap(driver))
}

func AddDriver(c fiber.Ctx) error {
	driver := models.Driver{}

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

	if tx := database.DB.Db.Create(&driver); tx.Error != nil {
		return c.Status(http.StatusInternalServerError).JSON(helpers.BuildError(tx.Error))
	}

	return c.Status(http.StatusCreated).JSON(helpers.ParseResultToMap(driver))
}

func UpdateDriver(c fiber.Ctx) error {
	driver := models.Driver{}

	if err := json.Unmarshal(c.Body(), &driver); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err.Error())
	}

	if driver.ID == 0 {
		return c.Status(http.StatusBadRequest).JSON(helpers.BuildError(fmt.Errorf("id is required")))
	}

	if tx := database.DB.Db.Model(&driver).Clauses(clause.Returning{}).Updates(driver); tx.Error != nil {
		return c.Status(http.StatusInternalServerError).JSON(helpers.BuildError(tx.Error))
	}

	return c.Status(http.StatusOK).JSON(helpers.ParseResultToMap(driver))
}

func DeleteDriver(c fiber.Ctx) error {
	id := c.Params("id")

	if tx := database.DB.Db.Delete(models.Driver{}, id); tx.Error != nil {
		return c.Status(http.StatusInternalServerError).JSON(helpers.BuildError(tx.Error))
	}

	return c.Status(http.StatusOK).JSON(helpers.ParseResultToMap(""))
}
