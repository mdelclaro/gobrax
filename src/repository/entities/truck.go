package entities

import (
	"github.com/shopspring/decimal"
)

type Truck struct {
	GormModel

	LicensePlate     string          `json:"licensePlate" validate:"required" gorm:"unique"`
	FuelUsed         decimal.Decimal `json:"fuelUsed"`
	DistanceTraveled decimal.Decimal `json:"distanceTraveled"`

	DriverID *int32  `json:"driverId" gorm:"unique"`
	Driver   *Driver `json:"driver" gorm:"foreignKey:DriverID"`
}
