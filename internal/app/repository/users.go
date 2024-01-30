package repository

import (
	"R_I_P_labs/internal/app/ds"
	"errors"
	"gorm.io/gorm"
)

func (r *Repository) AddUser(user *ds.User) error {
	return r.db.Create(user).Error
}
func (r *Repository) SaveUser(user *ds.User) error {
	return r.db.Save(user).Error
	
}

func (r *Repository) GetUserByLogin(login string) (*ds.User, error) {
	user := &ds.User{}
	if err := r.db.Where("login = ?", login).First(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (r *Repository) GetUserById(uuid string) (*ds.User, error) {
	user := &ds.User{}
	if err := r.db.Where("uuid = ?", uuid).First(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}
