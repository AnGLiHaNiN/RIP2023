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

type ChangeUserReq struct {
	Password *string `json:"password" binding:"omitempty,max=30"`
	Name     *string `json:"name" binding:"omitempty,max=60"`
	Email    *string `json:"email" binding:"omitempty,max=40"`
}

type ChangeComponentRequest struct {
	ComponentId string                `uri:"component_id" binding:"required,uuid"`
	Name        *string               `form:"name" json:"name" binding:"omitempty,max=100"`
	WorldName   *string               `form:"world_name" json:"world_name" binding:"omitempty,max=75"`
	Amount      *int                  `form:"amount" json:"amount"`
	Image       *multipart.FileHeader `form:"image" json:"image"`
	Properties  *string               `form:"properties" json:"properties" binding:"omitempty,max=200"`
}

type AddToMedicineRequest struct {
	ComponentId string `uri:"component_id" binding:"required,uuid"`
}

type GetAllMedicinesRequst struct {
	FormationDateStart *time.Time `form:"formation_date_start" json:"formation_date_start" time_format:"2006-01-02"`
	FormationDateEnd   *time.Time `form:"formation_date_end" json:"formation_date_end" time_format:"2006-01-02"`
	Status             string     `form:"status"`
}

type MedicineRequest struct {
	MedicineId string `uri:"medicine_id" binding:"required,uuid"`
}

type UpdateMedicineRequest struct {
	Name string `form:"name" json:"name" binding:"required,max=50"`
}

type ChangeCountReq struct {
	URI struct {
		ComponentId string `uri:"component_id" binding:"required,uuid"`
	}
	Count int `json:"count" binding:"required"`
}

type DeleteFromMedicineRequest struct {
	ComponentId string `uri:"component_id" binding:"required,uuid"`
}

type UserConfirmRequest struct {
	Confirm bool `form:"confirm" binding:"required"`
}

type ModeratorConfirmRequest struct {
	URI struct {
		MedicineId string `uri:"medicine_id" binding:"required,uuid"`
	}
	Confirm *bool `form:"confirm" binding:"required"`
}

type LoginReq struct {
	Login    string `form:"login" binding:"required,max=30"`
	Password string `form:"password" binding:"required,max=30"`
}

type RegisterReq struct {
	Login    string `form:"login" binding:"required,max=30"`
	Password string `form:"password" binding:"required,max=30"`
}

type VerificationReq struct {
	URI struct {
		MedicineId string `uri:"medicine_id" binding:"required,uuid"`
	}
	VerificationStatus *bool  `json:"verification_status" form:"verification_status" binding:"required"`
	Token              string `json:"token" form:"token" binding:"required"`
}
