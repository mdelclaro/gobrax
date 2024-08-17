package repository

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	db           *gorm.DB
	defaultJoins []string
}

func NewRepository(db *gorm.DB, joins ...string) *Repository {
	return &Repository{
		db:           db,
		defaultJoins: joins,
	}
}

func (r *Repository) Create(target any) error {
	res := r.db.Create(target)
	return r.HandleError(res)
}

func (r *Repository) FindById(target any, id int32, preloads ...string) error {
	res := r.DBWithPreloads(preloads).First(target, id)
	return r.HandleError(res)
}

func (r *Repository) FindAll(target any, preloads ...string) error {
	res := r.DBWithPreloads(preloads).Find(target)
	return r.HandleError(res)
}

func (r *Repository) Update(target any) error {
	res := r.db.
		Model(target).
		Clauses(clause.Returning{}).
		Updates(target)

	if res.RowsAffected == 0 {
		res.Error = fmt.Errorf("record not found")
	}

	return r.HandleError(res)
}

func (r *Repository) UpdateColumn(target any, column string, value any) error {
	res := r.db.
		Model(target).
		Clauses(clause.Returning{}).
		Update(column, value)

	return r.HandleError(res)
}

func (r *Repository) Delete(target any, id int32) error {
	res := r.db.Delete(target, id)
	if res.RowsAffected == 0 {
		res.Error = fmt.Errorf("record not found")
	}

	return r.HandleError(res)
}

func (r *Repository) HandleError(res *gorm.DB) error {
	if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
		err := fmt.Errorf("%w", res.Error)
		return err
	}

	return nil
}

func (r *Repository) DBWithPreloads(preloads []string) *gorm.DB {
	dbConn := r.db

	for _, join := range r.defaultJoins {
		dbConn = dbConn.Joins(join)
	}

	for _, preload := range preloads {
		dbConn = dbConn.Preload(preload)
	}

	return dbConn
}
