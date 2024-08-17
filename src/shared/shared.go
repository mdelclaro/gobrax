package shared

import (
	database "github.com/mdelclaro/gobrax/src/db"
	"github.com/mdelclaro/gobrax/src/repository"
	"github.com/mdelclaro/gobrax/src/repository/interfaces"
)

func InitRepo(joins ...string) interfaces.IRepository {
	return repository.NewRepository(database.StartDb().Db, joins...)
}
