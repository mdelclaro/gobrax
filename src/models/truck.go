package models

import (
	"github.com/shopspring/decimal"
)

type Truck struct {
	GormModel

	LicensePlate     string          `json:"licensePlate" validate:"required" gorm:"unique"`
	FuelUsed         decimal.Decimal `json:"fuelUsed" validate:"required"`
	DistanceTraveled decimal.Decimal `json:"distanceTraveled" validate:"required"`

	DriverID int32  `json:"driverId" validate:"required"`
	Driver   Driver `json:"driver" gorm:"foreignKey:DriverID"`
}
