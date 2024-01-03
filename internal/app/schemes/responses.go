package schemes

import (
	"lab3/internal/app/ds"
)

type AllComponentsResponse struct {
	Components []ds.Component `json:"components"`
}

type MedicineShort struct {
	UUID           string `json:"uuid"`
	ComponentCount int    `json:"component_count"`
}

type GetAllComponentsResponse struct {
	DraftMedicine *MedicineShort         `json:"draft_medicine"`
	Components            []ds.Component `json:"components"`
}

type AllMedicinesResponse struct {
	Medicines []MedicineOutput `json:"medicines"`
}

type MedicineResponse struct {
	Medicine MedicineOutput `json:"medicine"`
	Components    []ds.Component  `json:"components"`
}

type UpdateMedicineResponse struct {
	Medicine MedicineOutput  `json:"medicines"`
}

type MedicineOutput struct {
	UUID           string  `json:"uuid"`
	Status         string  `json:"status"`
	CreationDate   string  `json:"creation_date"`
	FormationDate  *string `json:"formation_date"`
	CompletionDate *string `json:"completion_date"`
	Moderator      *string `json:"moderator"`
	Customer       string  `json:"customer"`
	MedicineName      string  `json:"medicine_type"`
}

func ConvertMedicine(medicine *ds.Medicine) MedicineOutput {
	output := MedicineOutput{
		UUID:         medicine.UUID,
		Status:       medicine.Status,
		CreationDate: medicine.CreationDate.Format("2006-01-02 15:04:05"),
		MedicineName:    medicine.MedicineName,
		Customer:     medicine.Customer.Name,
	}

	if medicine.FormationDate != nil {
		formationDate := medicine.FormationDate.Format("2006-01-02 15:04:05")
		output.FormationDate = &formationDate
	}

	if medicine.CompletionDate != nil {
		completionDate := medicine.CompletionDate.Format("2006-01-02 15:04:05")
		output.CompletionDate = &completionDate
	}

	if medicine.Moderator != nil {
		output.Moderator = &medicine.Moderator.Name
	}

	return output
}