package repository

import (
	"errors"
	"strings"

	"gorm.io/gorm"

	"R_I_P_labs/internal/app/ds"
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

func (r *Repository) GetMedicineByName(FIO string) ([]ds.Medicine, error) {
	var medicines []ds.Medicine

	err := r.db.
		Where("LOWER(fio) LIKE ?", "%"+strings.ToLower(FIO)+"%").Where("is_deleted = ?", false).
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
	NotContent := ds.MedicineProduction{ComponentId: componentId, MedicineId: medicineId}
	err := r.db.Create(&NotContent).Error
	if err != nil {
		return err
	}
	return nil
}
