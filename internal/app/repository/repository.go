package repository

import (
	"gorm.io/gorm"
	"pharmaBlend/internal/app/ds"
)

func New(db *gorm.DB) *Repository {
	repo := &Repository{db}
	return repo
}

type Repository struct {
	db *gorm.DB
}

func (er *Repository) GetAllIngredients() ([]ds.PharmIngredient, error) {
	var equipments = make([]ds.PharmIngredient, 0)
	db := er.db.Find(&equipments, "Status = ?", "active")
	return equipments, db.Error
}

func (er *Repository) GetByID(id int) (*ds.PharmIngredient, error) {
	var equipment = &ds.PharmIngredient{}
	db := er.db.First(equipment, id).First(equipment, "Status = ?", "active")
	return equipment, db.Error
}

func (er *Repository) SearchIngredientsByTitle(title string) ([]ds.PharmIngredient, error) {
	var ingredients = make([]ds.PharmIngredient, 0)
	db := er.db.Find(&ingredients, "Status = ?", "active").Find(&ingredients, "name LIKE ?", "%"+title+"%")
	return ingredients, db.Error
}

func (er *Repository) DeleteIngredientsByID(id int) error {
	db := er.db.Exec("UPDATE pharm_ingredients SET status='delete' WHERE id = ?;", id)
	return db.Error
}
