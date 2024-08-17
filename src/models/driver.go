package models

type Driver struct {
	GormModel

	Name          string `json:"name" validate:"required"`
	LicenseNumber string `json:"licenseNumber" validate:"required" gorm:"unique"`
}
