package schemes

import (
	"R_I_P_labs/internal/app/ds"
	"R_I_P_labs/internal/app/repository"
	"time"
)

type AllComponentsResponse struct {
	Components []ds.Component `json:"components"`
}

type GetAllComponentsResponse struct {
	DraftMedicine *string        `json:"draft_medicine"`
	Components    []ds.Component `json:"components"`
}

type AllMedicinesResponse struct {
	Medicines []MedicineOutput `json:"medicines"`
}

type MedicineResponse struct {
	Medicine   MedicineOutput                  `json:"medicine"`
	Components []repository.ComponentWithCount `json:"components"`
}

type UpdateMedicineResponse struct {
	Medicine MedicineOutput `json:"medicines"`
}

type MedicineOutput struct {
	UUID               string  `json:"uuid"`
	Status             string  `json:"status"`
	CreationDate       string  `json:"creation_date"`
	FormationDate      *string `json:"formation_date"`
	CompletionDate     *string `json:"completion_date"`
	Moderator          *string `json:"moderator"`
	Customer           string  `json:"customer"`
	Name               *string `json:"name"`
	VerificationStatus *string `json:"verification_status"`
}

func ConvertMedicine(medicine *ds.Medicine) MedicineOutput {
	output := MedicineOutput{
		UUID:               medicine.UUID,
		Status:             medicine.Status,
		CreationDate:       medicine.CreationDate.Format(time.RFC3339),
		Name:               medicine.Name,
		Customer:           medicine.Customer.Login,
		VerificationStatus: medicine.VerificationStatus,
	}

	if medicine.FormationDate != nil {
		formationDate := medicine.FormationDate.Format(time.RFC3339)
		output.FormationDate = &formationDate
	}

	if medicine.CompletionDate != nil {
		completionDate := medicine.CompletionDate.Format(time.RFC3339)
		output.CompletionDate = &completionDate
	}

	if medicine.Moderator != nil {
		output.Moderator = &medicine.Moderator.Login
	}

	return output
}

type AddToMedicineResp struct {
	ComponentsCount int64 `json:"component_count"`
}

type AuthResp struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}
