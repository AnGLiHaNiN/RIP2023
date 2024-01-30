package repository

import (
	"errors"
	"strings"

	"gorm.io/gorm"

	"R_I_P_labs/internal/app/ds"
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

func (r *Repository) GetComponentByName(name string) ([]ds.Component, error) {
	var components []ds.Component

	err := r.db.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(name)+"%").
		Where("is_deleted = ?", false).
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
	existingRecord := ds.MedicineProduction{}
	err := r.db.Where(&ds.MedicineProduction{MedicineId: medicineId, ComponentId: componentId}).First(&existingRecord).Error

	if err == nil {
		existingRecord.Count++
		return r.db.Save(&existingRecord).Error
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		newRecord := ds.MedicineProduction{
			MedicineId:  medicineId,
			ComponentId: componentId,
		}
		return r.db.Create(&newRecord).Error
	}
	return err
}

func (r *Repository) ChangeCount(medicineId, componentId string, count int) error {
	condition := ds.MedicineProduction{
		MedicineId:  medicineId,
		ComponentId: componentId,
	}
	updateData := ds.MedicineProduction{
		Count: count,
	}

	return r.db.Model(&ds.MedicineProduction{}).Where(&condition).Updates(&updateData).Error
}
