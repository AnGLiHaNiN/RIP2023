package schemes

import (
	"lab3/internal/app/ds"
)

type AllMedicinesResponse struct {
	Medicines []ds.Medicine `json:"medicines"`
}

type ComponentShort struct {
	UUID           string `json:"uuid"`
	MedicineCount int    `json:"medicine_count"`
}

type GetAllMedicinesResponse struct {
	DraftComponent *ComponentShort         `json:"draft_component"`
	Medicines            []ds.Medicine `json:"medicines"`
}

type AllComponentsResponse struct {
	Components []ComponentOutput `json:"components"`
}

type ComponentResponse struct {
	Component ComponentOutput `json:"component"`
	Medicines    []ds.Medicine  `json:"medicines"`
}

type UpdateComponentResponse struct {
	Component ComponentOutput  `json:"components"`
}

type ComponentOutput struct {
	UUID           string  `json:"uuid"`
	Status         string  `json:"status"`
	CreationDate   string  `json:"creation_date"`
	FormationDate  *string `json:"formation_date"`
	CompletionDate *string `json:"completion_date"`
	Moderator      *string `json:"moderator"`
	Customer       string  `json:"customer"`
	ComponentName  string  `json:"component_name"`
}

func ConvertComponent(component *ds.Component) ComponentOutput {
	output := ComponentOutput{
		UUID:         component.UUID,
		Status:       component.Status,
		CreationDate: component.CreationDate.Format("2006-01-02 15:04:05"),
		ComponentName:    component.ComponentName,
		Customer:     component.Customer.Name,
	}

	if component.FormationDate != nil {
		formationDate := component.FormationDate.Format("2006-01-02 15:04:05")
		output.FormationDate = &formationDate
	}

	if component.CompletionDate != nil {
		completionDate := component.CompletionDate.Format("2006-01-02 15:04:05")
		output.CompletionDate = &completionDate
	}

	if component.Moderator != nil {
		output.Moderator = &component.Moderator.Name
	}

	return output
}