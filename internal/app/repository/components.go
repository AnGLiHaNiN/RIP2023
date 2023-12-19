package repository

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"

	"lab3/internal/app/ds"
)

func (r *Repository) GetAllComponents(formationDateStart, formationDateEnd *time.Time, status string) ([]ds.Component, error) {
	var components []ds.Component
	query := r.db.Preload("Customer").Preload("Moderator").
		Where("LOWER(status) LIKE ?", "%"+strings.ToLower(status)+"%").
		Where("status != ?", ds.DELETED)

	if formationDateStart != nil && formationDateEnd != nil {
		query = query.Where("formation_date BETWEEN ? AND ?", *formationDateStart, *formationDateEnd)
	} else if formationDateStart != nil {
		query = query.Where("formation_date >= ?", *formationDateStart)
	} else if formationDateEnd != nil {
		query = query.Where("formation_date <= ?", *formationDateEnd)
	}
	if err := query.Find(&components).Error; err != nil {
		return nil, err
	}
	return components, nil
}

func (r *Repository) GetDraftComponent(customerId string) (*ds.Component, error) {
	component := &ds.Component{}
	err := r.db.First(component, ds.Component{Status: ds.DRAFT, CustomerId: customerId}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return component, nil
}

func (r *Repository) CreateDraftComponent(customerId string) (*ds.Component, error) {
	component := &ds.Component{CreationDate: time.Now(), CustomerId: customerId, Status: ds.DRAFT}
	err := r.db.Create(component).Error
	if err != nil {
		return nil, err
	}
	return component, nil
}

func (r *Repository) GetComponentById(componentId, customerId string) (*ds.Component, error) {
	component := &ds.Component{}
	err := r.db.Preload("Moderator").Preload("Customer").
		Where("status != ?", ds.DELETED).
		First(component, ds.Component{UUID: componentId, CustomerId: customerId}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return component, nil
}

func (r *Repository) GetComponentContent(componentId string) ([]ds.Medicine, error) {
	var medicines []ds.Medicine

	err := r.db.Table("medicine_productions").
		Select("medicines.*").
		Joins("JOIN medicines ON medicine_productions.medicine_id = medicines.uuid").
		Where(ds.MedicineProduction{ComponentId: componentId}).
		Scan(&medicines).Error

	if err != nil {
		return nil, err
	}
	return medicines, nil
}

func (r *Repository) SaveComponent(component *ds.Component) error {
	err := r.db.Save(component).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) DeleteFromComponent(componentId, medicineId string) error {
	err := r.db.Delete(&ds.MedicineProduction{ComponentId: componentId, MedicineId: medicineId}).Error
	if err != nil {
		return err
	}
	return nil
}
