package repository

import (
	"errors"
	"strings"

	"gorm.io/gorm"

	"lab3/internal/app/ds"
)

func (r *Repository) GetComponentByID(id string) (*ds.Component, error) {
	component := &ds.Component{UUID: id}
	err := r.db.First(component, "is_deleted = ?", false).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return component, nil
}

func (r *Repository) AddComponent(component *ds.Component) error {
	err := r.db.Create(&component).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetComponentByName(Name string) ([]ds.Component, error) {
	var components []ds.Component

	err := r.db.
		Where("LOWER(name) LIKE ?", "%"+strings.ToLower(Name)+"%").Where("is_deleted = ?", false).
		Find(&components).Error

	if err != nil {
		return nil, err
	}

	return components, nil
}

func (r *Repository) SaveComponent(component *ds.Component) error {
	err := r.db.Save(component).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) AddToMedicine(medicineId, componentId string) error {
	MedProduction := ds.MedicineProduction{MedicineId: medicineId, ComponentId: componentId}
	err := r.db.Create(&MedProduction).Error
	if err != nil {
		return err
	}
	return nil
}
