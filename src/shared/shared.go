package shared

import (
	"github.com/mdelclaro/gobrax/src/repository"
	"github.com/mdelclaro/gobrax/src/repository/interfaces"
	"gorm.io/gorm"
)

func InitRepo(db *gorm.DB, joins ...string) interfaces.IRepository {
	return repository.NewRepository(db, joins...)
}
