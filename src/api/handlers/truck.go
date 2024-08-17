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

func GetAllTrucks(c fiber.Ctx) error {
	truck := models.Truck{}

	if tx := database.DB.Db.Find(&truck); tx.Error != nil {
		return c.Status(http.StatusInternalServerError).JSON(helpers.BuildError(tx.Error))
	}

	if truck.ID == 0 {
		return c.Status(http.StatusNoContent).JSON(helpers.ParseResultToMap(""))
	}

	return c.Status(http.StatusOK).JSON(helpers.ParseResultToMap(truck))
}

func GetTruckByID(c fiber.Ctx) error {
	truck := models.Truck{}
	id := c.Params("id")

	if tx := database.DB.Db.First(&truck, id); tx.Error != nil {
		return c.Status(http.StatusInternalServerError).JSON(helpers.BuildError(tx.Error))
	}

	if truck.ID == 0 {
		return c.Status(http.StatusNoContent).JSON(helpers.ParseResultToMap(""))
	}

	return c.Status(http.StatusOK).JSON(helpers.ParseResultToMap(truck))
}

func AddTruck(c fiber.Ctx) error {
	truck := models.Truck{}

	if err := json.Unmarshal(c.Body(), &truck); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(helpers.BuildError(err))
	}

	if err := helpers.Validator.Struct(truck); err != nil {
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

	if tx := database.DB.Db.Create(&truck); tx.Error != nil {
		return c.Status(http.StatusInternalServerError).JSON(helpers.BuildError(tx.Error))
	}

	return c.Status(http.StatusCreated).JSON(helpers.ParseResultToMap(truck))
}

func UpdateTruck(c fiber.Ctx) error {
	truck := models.Truck{}

	if err := json.Unmarshal(c.Body(), &truck); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err.Error())
	}

	if truck.ID == 0 {
		return c.Status(http.StatusBadRequest).JSON(helpers.BuildError(fmt.Errorf("id is required")))
	}

	if truck.DriverID != 0 {
		t := models.Truck{}

		if tx := database.DB.Db.First(&t, truck.ID); tx.Error != nil {
			return c.Status(http.StatusInternalServerError).JSON(helpers.BuildError(tx.Error))
		}

		if truck.ID == 0 {
			return c.Status(http.StatusInternalServerError).JSON(helpers.BuildError(fmt.Errorf("invalid truck id")))
		}

		if t.DriverID != 0 {
			return c.Status(http.StatusBadRequest).JSON(helpers.BuildError(fmt.Errorf("truck already in use by a driver")))
		}
	}

	if tx := database.DB.Db.Model(&truck).Clauses(clause.Returning{}).Updates(truck); tx.Error != nil {
		return c.Status(http.StatusInternalServerError).JSON(helpers.BuildError(tx.Error))
	}

	return c.Status(http.StatusOK).JSON(helpers.ParseResultToMap(truck))
}

func DeleteTruck(c fiber.Ctx) error {
	id := c.Params("id")

	if tx := database.DB.Db.Delete(models.Truck{}, id); tx.Error != nil {
		return c.Status(http.StatusInternalServerError).JSON(helpers.BuildError(tx.Error))
	}

	return c.Status(http.StatusOK).JSON(helpers.ParseResultToMap(""))
}
