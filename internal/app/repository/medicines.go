package repository

import (
	"errors"
	"strings"

	"gorm.io/gorm"

	"lab3/internal/app/ds"
)

func (r *Repository) GetMedicineByID(id string) (*ds.Medicine, error) {
	medicine := &ds.Medicine{UUID: id}
	err := r.db.First(medicine, "is_deleted = ?", false).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return medicine, nil
}

func (r *Repository) AddMedicine(medicine *ds.Medicine) error {
	err := r.db.Create(&medicine).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetMedicineByName(Name string) ([]ds.Medicine, error) {
	var medicines []ds.Medicine

	err := r.db.
		Where("LOWER(name) LIKE ?", "%"+strings.ToLower(Name)+"%").Where("is_deleted = ?", false).
		Find(&medicines).Error

	if err != nil {
		return nil, err
	}

	return medicines, nil
}

func (r *Repository) SaveMedicine(medicine *ds.Medicine) error {
	err := r.db.Save(medicine).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) AddToComponent(componentId, medicineId string) error {
	MedProd := ds.MedicineProduction{ComponentId: componentId, MedicineId: medicineId}
	err := r.db.Create(&MedProd).Error
	if err != nil {
		return err
	}
	return nil
}
