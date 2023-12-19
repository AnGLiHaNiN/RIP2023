package schemes

import (
	"R_I_P_labs/internal/app/ds"

	"mime/multipart"
	"time"
)

type MedicineRequest struct {
	MedicineId string `uri:"medicine_id" binding:"required,uuid"`
}

type GetAllMedicinesRequest struct {
	Name string `form:"name"`
}

type AddMedicineRequest struct {
	ds.Medicine
	Image *multipart.FileHeader `form:"image" json:"image"`
}

type ChangeMedicineRequest struct {
	MedicineId string                `uri:"medicine_id" binding:"required,uuid"`
	Name         *string               `form:"name" json:"name" binding:"omitempty,max=100"`
	Manufacturer       *string               `form:"manufacturer" json:"manufacturer" binding:"omitempty,max=75"`
	Amount         *int                  `form:"amount" json:"amount"`
	Image       *multipart.FileHeader `form:"image" json:"image"`
	Dosage      *string               `form:"dosage" json:"dosage" binding:"omitempty,max=100"`
}

type AddToComponentRequest struct {
	MedicineId string `uri:"medicine_id" binding:"required,uuid"`
}

type GetAllComponentsRequst struct {
	FormationDateStart *time.Time `form:"formation_date_start" json:"formation_date_start" time_format:"2006-01-02 15:04:05"`
	FormationDateEnd   *time.Time `form:"formation_date_end" json:"formation_date_end" time_format:"2006-01-02 15:04:05"`
	Status             string     `form:"status"`
}

type ComponentRequest struct {
	ComponentId string `uri:"component_id" binding:"required,uuid"`
}

type UpdateComponentRequest struct {
	URI struct {
		ComponentId string `uri:"component_id" binding:"required,uuid"`
	}
	ComponentName string `form:"component_name" json:"component_name" binding:"required,max=50"`
}

type DeleteFromComponentRequest struct {
	ComponentId string `uri:"component_id" binding:"required,uuid"`
	MedicineId    string `uri:"medicine_id" binding:"required,uuid"`
}

type UserConfirmRequest struct {
	Confirm bool `form:"confirm" binding:"required"`
}

type ModeratorConfirmRequest struct {
	URI struct {
		ComponentId string `uri:"component_id" binding:"required,uuid"`
	}
	Confirm bool `form:"confirm" binding:"required"`
}

type LoginReq struct {
	Login    string `form:"login" binding:"required,max=30"`
	Password string `form:"password" binding:"required,max=30"`
}

type RegisterReq struct {
	Login    string `form:"login" binding:"required,max=30"`
	Password string `form:"password" binding:"required,max=30"`
}