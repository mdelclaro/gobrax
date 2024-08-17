package interfaces

import "gorm.io/gorm"

type IRepository interface {
	Create(target any) error
	FindById(target any, id int32, preloads ...string) error
	FindAll(target any, preloads ...string) error
	Update(target any) error
	UpdateColumn(target any, column string, value any) error
	Delete(target any, id int32) error
	HandleError(res *gorm.DB) error
	DBWithPreloads(preloads []string) *gorm.DB
}
