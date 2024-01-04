package schemes

import (
	"R_I_P_labs/internal/app/ds"

	"mime/multipart"
	"time"
)

type ComponentRequest struct {
	ComponentId string `uri:"component_id" binding:"required,uuid"`
}

type GetAllComponentsRequest struct {
	Name string `form:"name"`
}

type AddComponentRequest struct {
	ds.Component
	Image *multipart.FileHeader `form:"image" json:"image"`
}

type ChangeComponentRequest struct {
	ComponentId string                `uri:"component_id" binding:"required,uuid"`
	Name         *string               `form:"name" json:"name" binding:"omitempty,max=100"`
	WorldName       *string               `form:"world_name" json:"world_name" binding:"omitempty,max=75"`
	Amount         *int                  `form:"amount" json:"amount"`
	Image       *multipart.FileHeader `form:"image" json:"image"`
	Properties      *string               `form:"properties" json:"properties" binding:"omitempty,max=100"`
}

type AddToMedicineRequest struct {
	ComponentId string `uri:"component_id" binding:"required,uuid"`
}

type GetAllMedicinesRequst struct {
	FormationDateStart *time.Time `form:"formation_date_start" json:"formation_date_start" time_format:"2006-01-02 15:04:05"`
	FormationDateEnd   *time.Time `form:"formation_date_end" json:"formation_date_end" time_format:"2006-01-02 15:04:05"`
	Status             string     `form:"status"`
}

type MedicineRequest struct {
	MedicineId string `uri:"medicine_id" binding:"required,uuid"`
}

type UpdateMedicineRequest struct {
	URI struct {
		MedicineId string `uri:"medicine_id" binding:"required,uuid"`
	}
	MedicineName string `form:"medicine_name" json:"medicine_name" binding:"required,max=50"`
}

type DeleteFromMedicineRequest struct {
	MedicineId string `uri:"medicine_id" binding:"required,uuid"`
	ComponentId    string `uri:"component_id" binding:"required,uuid"`
}

type UserConfirmRequest struct {
	Confirm bool `form:"confirm" binding:"required"`
}

type ModeratorConfirmRequest struct {
	URI struct {
		MedicineId string `uri:"medicine_id" binding:"required,uuid"`
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