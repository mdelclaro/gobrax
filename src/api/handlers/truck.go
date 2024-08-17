package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v3"
	"github.com/mdelclaro/gobrax/src/api/helpers"
	"github.com/mdelclaro/gobrax/src/repository/entities"
	"github.com/mdelclaro/gobrax/src/shared"
)

func GetAllTrucks(c fiber.Ctx) error {
	trucks := []entities.Truck{}

	if err := shared.InitRepo("Driver").FindAll(&trucks, "Driver"); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(helpers.BuildError(err))
	}

	if len(trucks) == 0 {
		return c.Status(http.StatusNoContent).JSON(helpers.ParseResultToMap(""))
	}

	return c.Status(http.StatusOK).JSON(helpers.ParseResultToMap(trucks))
}

func GetTruckByID(c fiber.Ctx) error {
	truck := entities.Truck{}

	id := c.Params("id")
	parsedId, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(helpers.BuildError(fmt.Errorf("invalid id provided: %s", err.Error())))
	}

	if err := shared.InitRepo("Driver").FindById(&truck, int32(parsedId), "Driver"); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(helpers.BuildError(err))
	}

	if truck.ID == 0 {
		return c.Status(http.StatusNoContent).JSON(helpers.ParseResultToMap(""))
	}

	return c.Status(http.StatusOK).JSON(helpers.ParseResultToMap(truck))
}

func AddTruck(c fiber.Ctx) error {
	truck := entities.Truck{}

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

	if err := shared.InitRepo().Create(&truck); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(helpers.BuildError(err))
	}

	return c.Status(http.StatusCreated).JSON(helpers.ParseResultToMap(truck))
}

func UpdateTruck(c fiber.Ctx) error {
	truck := entities.Truck{}

	if err := json.Unmarshal(c.Body(), &truck); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err.Error())
	}

	if truck.ID == 0 {
		return c.Status(http.StatusBadRequest).JSON(helpers.BuildError(fmt.Errorf("id is required")))
	}

	if truck.DriverID != nil {
		return c.Status(http.StatusBadRequest).JSON(helpers.BuildError(fmt.Errorf("can't directly update driver id")))
	}

	if err := shared.InitRepo().Update(&truck); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(helpers.BuildError(err))
	}

	// return updated truck with driver association
	if err := shared.InitRepo("Driver").FindById(&truck, truck.ID, "Driver"); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(helpers.BuildError(err))
	}

	return c.Status(http.StatusOK).JSON(helpers.ParseResultToMap(truck))
}

func DeleteTruck(c fiber.Ctx) error {
	id := c.Params("id")
	parsedId, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(helpers.BuildError(fmt.Errorf("invalid id provided: %s", err.Error())))
	}

	if err := shared.InitRepo().Delete(entities.Truck{}, int32(parsedId)); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(helpers.BuildError(err))
	}

	return c.Status(http.StatusOK).JSON(helpers.ParseResultToMap(""))
}

func UpdateTruckDriver(c fiber.Ctx) error {
	truck := entities.Truck{}
	driver := entities.Driver{}

	truckId := c.Params("id")
	driverId := c.Query("driverId")

	parsedTruckId, err := strconv.Atoi(truckId)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(helpers.BuildError(fmt.Errorf("invalid truck id provided: %s", err.Error())))
	}

	if err := shared.InitRepo().FindById(&truck, int32(parsedTruckId)); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(helpers.BuildError(err))
	}

	if truck.ID == 0 {
		return c.Status(http.StatusNotFound).JSON(helpers.BuildError(fmt.Errorf("truck not found")))
	}

	parsedDriverId, err := strconv.Atoi(driverId)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(helpers.BuildError(fmt.Errorf("invalid driver id provided: %s", err.Error())))
	}

	if err := shared.InitRepo().FindById(&driver, int32(parsedDriverId)); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(helpers.BuildError(err))
	}

	if driver.ID == 0 || !driver.IsActive {
		return c.Status(http.StatusBadRequest).JSON(helpers.BuildError(fmt.Errorf("invalid driver provided")))
	}

	if err := shared.InitRepo().UpdateColumn(&truck, "driver_id", int32(parsedDriverId)); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(helpers.BuildError(err))
	}

	// return updated truck with driver association
	if err := shared.InitRepo("Driver").FindById(&truck, int32(parsedTruckId), "Driver"); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(helpers.BuildError(err))
	}

	return c.Status(http.StatusOK).JSON(helpers.ParseResultToMap(truck))
}
