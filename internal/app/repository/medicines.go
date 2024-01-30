package repository

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"

	"R_I_P_labs/internal/app/ds"
)

func (r *Repository) GetAllMedicines(customerId *string, formationDateStart, formationDateEnd *time.Time, status string) ([]ds.Medicine, error) {
	var medicines []ds.Medicine
	query := r.db.Preload("Customer").Preload("Moderator").
		Where("LOWER(status) LIKE ?", "%"+strings.ToLower(status)+"%").
		Where("status != ? AND status != ?", ds.DELETED, ds.DRAFT)

	if customerId != nil {
		query = query.Where("customer_id = ?", *customerId)
	}

	if formationDateStart != nil && formationDateEnd != nil {
		query = query.Where("formation_date BETWEEN ? AND ?", *formationDateStart, *formationDateEnd)
	} else if formationDateStart != nil {
		query = query.Where("formation_date >= ?", *formationDateStart)
	} else if formationDateEnd != nil {
		query = query.Where("formation_date <= ?", *formationDateEnd)
	}
	if err := query.Find(&medicines).Error; err != nil {
		return nil, err
	}
	return medicines, nil
}

func (r *Repository) GetDraftMedicine(customerId string) (*ds.Medicine, error) {
	medicine := &ds.Medicine{}
	err := r.db.First(medicine, ds.Medicine{Status: ds.DRAFT, CustomerId: customerId}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return medicine, nil
}

func (r *Repository) CreateDraftMedicine(customerId string) (*ds.Medicine, error) {
	medicine := &ds.Medicine{CreationDate: time.Now(), CustomerId: customerId, Status: ds.DRAFT}
	err := r.db.Create(medicine).Error
	if err != nil {
		return nil, err
	}
	return medicine, nil
}

func (r *Repository) GetMedicineById(medicineId string, userId *string) (*ds.Medicine, error) {
	medicine := &ds.Medicine{}
	query := r.db.Preload("Moderator").Preload("Customer").
		Where("status != ?", ds.DELETED)
	if userId != nil {
		query = query.Where("customer_id = ?", userId)
	}
	err := query.First(medicine, ds.Medicine{UUID: medicineId}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return medicine, nil
}

type ComponentWithCount struct {
	ds.Component
	Count int `json:"count"`
}

func (r *Repository) GetMedicineProduction(medicineId string) ([]ComponentWithCount, error) {
	var components []ComponentWithCount

	err := r.db.Table("medicine_productions").
		Select("components.*, medicine_productions.count as count").
		Joins("JOIN components ON medicine_productions.component_id = components.uuid").
		Where(ds.MedicineProduction{MedicineId: medicineId}).
		Scan(&components).Error

	if err != nil {
		return nil, err
	}
	return components, nil
}

func (r *Repository) SaveMedicine(medicine *ds.Medicine) error {
	return r.db.Save(medicine).Error
}

func (r *Repository) DeleteFromMedicine(medicineId, ComponentId string) error {
	return r.db.Delete(&ds.MedicineProduction{
		MedicineId:  medicineId,
		ComponentId: ComponentId,
	}).Error
}